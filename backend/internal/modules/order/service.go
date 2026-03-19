package order

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, req CreateOrderReq) (*OrderResp, error) {
	totalAmount := decimal.Zero
	items := make([]OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		unitPrice, err := decimal.NewFromString(item.UnitPrice)
		if err != nil {
			return nil, apperrors.ErrInvalidParams
		}
		itemTotal := unitPrice.Mul(decimal.NewFromInt(int64(item.Quantity)))
		totalAmount = totalAmount.Add(itemTotal)

		items = append(items, OrderItem{
			ItemType:   item.ItemType,
			ItemID:     item.ItemID,
			ItemName:   item.ItemName,
			SkuID:      item.SkuID,
			Quantity:   item.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: itemTotal,
		})
	}

	expiresAt := time.Now().Add(30 * time.Minute)

	order := &Order{
		OrderNo:        GenerateOrderNo(),
		UserID:         userID,
		Type:           req.Type,
		Status:         StatusPending,
		TotalAmount:    totalAmount,
		DiscountAmount: decimal.Zero,
		PayAmount:      totalAmount,
		BalancePaid:    decimal.Zero,
		WechatPaid:     decimal.Zero,
		Remark:         req.Remark,
		ExpiresAt:      &expiresAt,
		Items:          items,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToOrderResp(order)
	return &resp, nil
}

// CreateOrderDirect creates an order from pre-built model (used by booking service).
func (s *Service) CreateOrderDirect(ctx context.Context, order *Order) error {
	return s.repo.CreateOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, userID, orderID int64) (*OrderResp, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if order.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	resp := ToOrderResp(order)
	return &resp, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID int64) (*Order, error) {
	return s.repo.GetByID(ctx, orderID)
}

func (s *Service) ListUserOrders(ctx context.Context, userID int64, status string, offset, limit int) ([]OrderResp, int64, error) {
	orders, total, err := s.repo.ListByUserID(ctx, userID, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]OrderResp, 0, len(orders))
	for i := range orders {
		result = append(result, ToOrderResp(&orders[i]))
	}
	return result, total, nil
}

func (s *Service) ListAllOrders(ctx context.Context, status, orderType string, offset, limit int) ([]OrderResp, int64, error) {
	orders, total, err := s.repo.ListAll(ctx, status, orderType, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]OrderResp, 0, len(orders))
	for i := range orders {
		result = append(result, ToOrderResp(&orders[i]))
	}
	return result, total, nil
}

func (s *Service) CancelOrder(ctx context.Context, userID, orderID int64, reason string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return apperrors.ErrRecordNotFound
	}

	if order.UserID != userID {
		return apperrors.ErrForbidden
	}

	if err := ValidateTransition(order.Status, StatusCancelled); err != nil {
		return apperrors.New(40001, err.Error())
	}

	return s.repo.UpdateStatus(ctx, orderID, order.Status, StatusCancelled)
}

func (s *Service) RequestRefund(ctx context.Context, userID, orderID int64, reason string) (*RefundResp, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if order.UserID != userID {
		return nil, apperrors.ErrForbidden
	}

	if err := ValidateTransition(order.Status, StatusRefunding); err != nil {
		return nil, apperrors.New(40001, err.Error())
	}

	// Update order status to refunding
	if err := s.repo.UpdateStatus(ctx, orderID, order.Status, StatusRefunding); err != nil {
		return nil, apperrors.ErrInternal
	}

	refund := &Refund{
		RefundNo:      GenerateRefundNo(),
		OrderID:       orderID,
		UserID:        userID,
		Amount:        order.PayAmount,
		BalanceRefund: order.BalancePaid,
		WechatRefund:  order.WechatPaid,
		Reason:        reason,
		Status:        "pending",
	}

	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToRefundResp(refund)
	return &resp, nil
}

func (s *Service) ReviewRefund(ctx context.Context, adminID, refundID int64, action, remark string) (*RefundResp, error) {
	refund, err := s.repo.GetRefundByID(ctx, refundID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if refund.Status != "pending" {
		return nil, apperrors.New(40001, "退款单已处理")
	}

	now := time.Now()
	refund.ReviewedBy = &adminID
	refund.ReviewedAt = &now

	if action == "approve" {
		refund.Status = "approved"
		// Update order status to refunded
		_ = s.repo.UpdateStatus(ctx, refund.OrderID, StatusRefunding, StatusRefunded)
	} else {
		refund.Status = "rejected"
		// Revert order status to paid
		_ = s.repo.UpdateStatus(ctx, refund.OrderID, StatusRefunding, StatusPaid)
	}

	if err := s.repo.UpdateRefund(ctx, refund); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToRefundResp(refund)
	return &resp, nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int64, fromStatus, toStatus string) error {
	return s.repo.UpdateStatus(ctx, orderID, fromStatus, toStatus)
}

func (s *Service) UpdateOrderPaid(ctx context.Context, orderID int64, balancePaid, wechatPaid decimal.Decimal) error {
	return s.repo.UpdateOrderPaid(ctx, orderID, balancePaid, wechatPaid)
}
