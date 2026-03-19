package payment

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	repo         *Repository
	wallet       WalletOperator
	wechat       WechatPayer
	orderUpdater OrderUpdater
}

func NewService(repo *Repository, wallet WalletOperator, wechat WechatPayer, orderUpdater OrderUpdater) *Service {
	return &Service{
		repo:         repo,
		wallet:       wallet,
		wechat:       wechat,
		orderUpdater: orderUpdater,
	}
}

// GeneratePaymentNo creates a unique payment number.
func GeneratePaymentNo() string {
	return "PAY" + time.Now().Format("20060102150405") + uuid.New().String()[:8]
}

// PreparePayment handles pre-payment: calculates amounts, picks method, creates payment record.
func (s *Service) PreparePayment(ctx context.Context, userID int64, req PreparePayReq, openID string) (*PreparePayResp, error) {
	// Get order info (via order updater - but we need amount from order)
	// For now, we'll compute from the payment request
	paymentNo := GeneratePaymentNo()

	// Calculate amounts based on method
	var balanceAmount, wechatAmount, totalAmount decimal.Decimal

	// We need the order amount - get from the existing payment or order
	// The handler will pass the order amount
	totalAmount, _ = decimal.NewFromString(req.BalanceAmount)
	if totalAmount.IsZero() {
		return nil, apperrors.ErrInvalidParams
	}

	switch req.Method {
	case MethodBalance:
		balance, err := s.wallet.GetBalance(ctx, userID)
		if err != nil {
			return nil, apperrors.ErrInternal
		}
		if balance.LessThan(totalAmount) {
			return nil, apperrors.ErrInsufficientBalance
		}
		balanceAmount = totalAmount
		wechatAmount = decimal.Zero

	case MethodWechat:
		balanceAmount = decimal.Zero
		wechatAmount = totalAmount

	case MethodCombo:
		balance, err := s.wallet.GetBalance(ctx, userID)
		if err != nil {
			return nil, apperrors.ErrInternal
		}
		if balance.GreaterThanOrEqual(totalAmount) {
			balanceAmount = totalAmount
			wechatAmount = decimal.Zero
		} else {
			balanceAmount = balance
			wechatAmount = totalAmount.Sub(balance)
		}
	}

	// Create payment record
	payment := &Payment{
		PaymentNo:     paymentNo,
		OrderID:       req.OrderID,
		UserID:        userID,
		Method:        req.Method,
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
		Method:        req.Method,
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

// HandleWechatCallback processes wechat payment callback (idempotent).
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
