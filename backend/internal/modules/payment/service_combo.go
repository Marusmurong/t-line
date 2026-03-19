package payment

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// WalletOperator abstracts wallet operations from auth module.
type WalletOperator interface {
	DeductBalance(ctx context.Context, userID int64, amount decimal.Decimal, orderID int64, remark string) error
	FreezeBalance(ctx context.Context, userID int64, amount decimal.Decimal) error
	UnfreezeBalance(ctx context.Context, userID int64, amount decimal.Decimal) error
	RefundBalance(ctx context.Context, userID int64, amount decimal.Decimal, orderID int64, remark string) error
	GetBalance(ctx context.Context, userID int64) (decimal.Decimal, error)
}

// WechatPayer abstracts wechat payment SDK.
type WechatPayer interface {
	CreatePrepayOrder(orderNo string, amount decimal.Decimal, description string, openID string) (*WechatPayParams, error)
	RefundOrder(paymentNo string, refundNo string, totalAmount, refundAmount decimal.Decimal) error
}

// OrderUpdater abstracts order status updates.
type OrderUpdater interface {
	UpdateOrderPaid(ctx context.Context, orderID int64, balancePaid, wechatPaid decimal.Decimal) error
	UpdateOrderStatus(ctx context.Context, orderID int64, fromStatus, toStatus string) error
}

// PayByBalance handles pure balance payment.
func (s *Service) PayByBalance(ctx context.Context, userID int64, orderID int64, amount decimal.Decimal, paymentNo string) error {
	// Deduct balance with optimistic lock
	if err := s.wallet.DeductBalance(ctx, userID, amount, orderID, "余额支付"); err != nil {
		return err
	}

	// Update payment record
	payment, err := s.repo.GetPaymentByNo(ctx, paymentNo)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePaymentStatus(ctx, payment.ID, PayStatusSuccess, ""); err != nil {
		return err
	}

	// Update order as paid
	return s.orderUpdater.UpdateOrderPaid(ctx, orderID, amount, decimal.Zero)
}

// PayByCombo handles combo payment: freeze balance -> create wechat prepay -> wait callback.
func (s *Service) PayByCombo(ctx context.Context, userID int64, orderID int64, balanceAmount, wechatAmount decimal.Decimal, paymentNo, orderNo, openID string) (*WechatPayParams, error) {
	// Freeze balance portion
	if balanceAmount.GreaterThan(decimal.Zero) {
		if err := s.wallet.FreezeBalance(ctx, userID, balanceAmount); err != nil {
			return nil, err
		}
	}

	// Create wechat prepay order
	params, err := s.wechat.CreatePrepayOrder(orderNo, wechatAmount, "T-Line 场馆预订", openID)
	if err != nil {
		// Unfreeze balance on failure
		if balanceAmount.GreaterThan(decimal.Zero) {
			_ = s.wallet.UnfreezeBalance(ctx, userID, balanceAmount)
		}
		return nil, err
	}

	return params, nil
}

// ConfirmComboPayment is called when wechat callback confirms payment.
func (s *Service) ConfirmComboPayment(ctx context.Context, payment *Payment) error {
	now := time.Now()
	payment.Status = PayStatusSuccess
	payment.PaidAt = &now

	return s.repo.db.Transaction(func(tx *gorm.DB) error {
		// Update payment
		if err := tx.WithContext(ctx).Save(payment).Error; err != nil {
			return err
		}

		// If balance was frozen, deduct it
		if payment.BalanceAmount.GreaterThan(decimal.Zero) {
			if err := s.wallet.DeductBalance(ctx, payment.UserID, payment.BalanceAmount, payment.OrderID, "组合支付-余额部分"); err != nil {
				return err
			}
			_ = s.wallet.UnfreezeBalance(ctx, payment.UserID, payment.BalanceAmount)
		}

		// Update order
		return s.orderUpdater.UpdateOrderPaid(ctx, payment.OrderID, payment.BalanceAmount, payment.WechatAmount)
	})
}

// RefundPayment handles refund: returns balance + wechat refund proportionally.
func (s *Service) RefundPayment(ctx context.Context, orderID int64, refundAmount decimal.Decimal) error {
	payment, err := s.repo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	// Calculate proportional refund
	totalPaid := payment.BalanceAmount.Add(payment.WechatAmount)
	if totalPaid.IsZero() {
		return nil
	}

	balanceRefund := decimal.Zero
	wechatRefund := decimal.Zero

	if payment.BalanceAmount.GreaterThan(decimal.Zero) {
		balanceRefund = refundAmount.Mul(payment.BalanceAmount).Div(totalPaid).Round(2)
	}
	wechatRefund = refundAmount.Sub(balanceRefund)

	// Refund balance portion
	if balanceRefund.GreaterThan(decimal.Zero) {
		if err := s.wallet.RefundBalance(ctx, payment.UserID, balanceRefund, orderID, "退款-余额部分"); err != nil {
			return err
		}
	}

	// Refund wechat portion
	if wechatRefund.GreaterThan(decimal.Zero) {
		refundNo := "RF" + payment.PaymentNo
		if err := s.wechat.RefundOrder(payment.PaymentNo, refundNo, payment.WechatAmount, wechatRefund); err != nil {
			return err
		}
	}

	// Update payment status
	return s.repo.UpdatePaymentStatus(ctx, payment.ID, PayStatusRefunded, "")
}
