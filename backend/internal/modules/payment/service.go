package payment

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"github.com/t-line/backend/internal/pkg/logger"
	"gorm.io/gorm"
)

// OrderQuerier abstracts reading order data (to avoid importing order module).
type OrderQuerier interface {
	GetOrderByID(ctx context.Context, orderID int64) (OrderInfo, error)
}

// OrderInfo is a minimal view of an order used by the payment module.
type OrderInfo struct {
	ID        int64
	UserID    int64
	Status    string
	PayAmount decimal.Decimal
	Type      string
}

type Service struct {
	repo         *Repository
	wallet       WalletOperator
	wechat       WechatPayer
	orderUpdater OrderUpdater
	orderQuerier OrderQuerier
}

func NewService(repo *Repository, wallet WalletOperator, wechat WechatPayer, orderUpdater OrderUpdater, orderQuerier OrderQuerier) *Service {
	return &Service{
		repo:         repo,
		wallet:       wallet,
		wechat:       wechat,
		orderUpdater: orderUpdater,
		orderQuerier: orderQuerier,
	}
}

// GeneratePaymentNo creates a unique payment number.
func GeneratePaymentNo() string {
	return "PAY" + time.Now().Format("20060102150405") + uuid.New().String()[:8]
}

// PreparePayment handles pre-payment: verifies order ownership, calculates amounts from server-side data.
func (s *Service) PreparePayment(ctx context.Context, userID int64, req PreparePayReq, openID string) (*PreparePayResp, error) {
	// 1. Query order from backend — never trust client amount
	orderInfo, err := s.orderQuerier.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	// 2. Verify order belongs to current user
	if orderInfo.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	// 3. Verify order is pending
	if orderInfo.Status != "pending" {
		return nil, apperrors.New(40001, "订单状态不允许支付")
	}

	totalAmount := orderInfo.PayAmount
	if totalAmount.LessThanOrEqual(decimal.Zero) {
		return nil, apperrors.ErrInvalidParams
	}

	paymentNo := GeneratePaymentNo()

	// 4. Query wallet balance to decide payment method automatically
	balance, err := s.wallet.GetBalance(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	var balanceAmount, wechatAmount decimal.Decimal
	var method string

	if balance.GreaterThanOrEqual(totalAmount) {
		// Pure balance payment
		method = MethodBalance
		balanceAmount = totalAmount
		wechatAmount = decimal.Zero
	} else if balance.GreaterThan(decimal.Zero) {
		// Combo payment
		method = MethodCombo
		balanceAmount = balance
		wechatAmount = totalAmount.Sub(balance)
	} else {
		// Pure wechat payment
		method = MethodWechat
		balanceAmount = decimal.Zero
		wechatAmount = totalAmount
	}

	// 5. Create payment record
	payment := &Payment{
		PaymentNo:     paymentNo,
		OrderID:       req.OrderID,
		UserID:        userID,
		Method:        method,
		Amount:        totalAmount,
		BalanceAmount: balanceAmount,
		WechatAmount:  wechatAmount,
		Status:        PayStatusPending,
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := &PreparePayResp{
		PaymentNo:     paymentNo,
		Method:        method,
		TotalAmount:   totalAmount.StringFixed(2),
		BalanceAmount: balanceAmount.StringFixed(2),
		WechatAmount:  wechatAmount.StringFixed(2),
		Paid:          false,
	}

	// Pure balance payment: execute immediately
	if wechatAmount.IsZero() && balanceAmount.GreaterThan(decimal.Zero) {
		if err := s.PayByBalance(ctx, userID, req.OrderID, balanceAmount, paymentNo); err != nil {
			return nil, apperrors.ErrPaymentFailed
		}
		resp.Paid = true
		return resp, nil
	}

	// Wechat or combo: create prepay order
	if wechatAmount.GreaterThan(decimal.Zero) {
		orderNo := paymentNo // Use payment no as trade no
		params, err := s.PayByCombo(ctx, userID, req.OrderID, balanceAmount, wechatAmount, paymentNo, orderNo, openID)
		if err != nil {
			return nil, apperrors.ErrPaymentFailed
		}
		resp.WechatParams = params
	}

	return resp, nil
}

// HandleWechatCallback processes wechat payment callback with signature verification and amount check.
func (s *Service) HandleWechatCallback(ctx context.Context, req WechatCallbackReq) error {
	payment, err := s.repo.GetPaymentByNo(ctx, req.OutTradeNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Idempotent: ignore unknown payments
		}
		return apperrors.ErrInternal
	}

	// Idempotent check: already processed
	if payment.Status == PayStatusSuccess {
		return nil
	}

	// Verify callback amount matches payment record (amount in cents)
	expectedCents := payment.WechatAmount.Mul(decimal.NewFromInt(100)).IntPart()
	if req.Amount != expectedCents {
		logger.L.Errorw("wechat callback amount mismatch",
			"payment_no", payment.PaymentNo,
			"expected_cents", expectedCents,
			"actual_cents", req.Amount,
		)
		return apperrors.New(40001, "回调金额与订单不匹配")
	}

	if req.TradeState != "SUCCESS" {
		_ = s.repo.UpdatePaymentStatus(ctx, payment.ID, PayStatusFailed, req.TransactionID)
		// Unfreeze balance if combo
		if payment.BalanceAmount.GreaterThan(decimal.Zero) {
			_ = s.wallet.UnfreezeBalance(ctx, payment.UserID, payment.BalanceAmount)
		}
		return nil
	}

	payment.WechatTradeNo = req.TransactionID

	return s.ConfirmComboPayment(ctx, payment)
}

// ListUserCoupons returns user's coupons.
func (s *Service) ListUserCoupons(ctx context.Context, userID int64) ([]CouponResp, error) {
	unused := CouponUnused
	userCoupons, err := s.repo.ListUserCoupons(ctx, userID, &unused)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	return toCouponRespList(userCoupons), nil
}

// ListAvailableCoupons returns coupons available for a specific order.
func (s *Service) ListAvailableCoupons(ctx context.Context, userID int64, orderType string, amount string) ([]CouponResp, error) {
	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	userCoupons, err := s.repo.ListAvailableCoupons(ctx, userID, orderType, amountDec)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	// Filter by applicable types
	var filtered []UserCoupon
	for _, uc := range userCoupons {
		if uc.Coupon == nil {
			continue
		}
		var types []string
		_ = json.Unmarshal(uc.Coupon.ApplicableTypes, &types)
		if len(types) == 0 {
			filtered = append(filtered, uc)
			continue
		}
		for _, t := range types {
			if t == orderType {
				filtered = append(filtered, uc)
				break
			}
		}
	}

	return toCouponRespList(filtered), nil
}

func toCouponRespList(userCoupons []UserCoupon) []CouponResp {
	result := make([]CouponResp, 0, len(userCoupons))
	for _, uc := range userCoupons {
		if uc.Coupon == nil {
			continue
		}
		var types []string
		_ = json.Unmarshal(uc.Coupon.ApplicableTypes, &types)

		result = append(result, CouponResp{
			ID:              uc.Coupon.ID,
			UserCouponID:    uc.ID,
			Name:            uc.Coupon.Name,
			Type:            uc.Coupon.Type,
			Value:           uc.Coupon.Value.StringFixed(2),
			MinAmount:       uc.Coupon.MinAmount.StringFixed(2),
			ApplicableTypes: types,
			Status:          uc.Status,
			ExpiresAt:       uc.ExpiresAt,
		})
	}
	return result
}
