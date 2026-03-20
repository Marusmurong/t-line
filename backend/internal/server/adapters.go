package server

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	authmod "github.com/t-line/backend/internal/modules/auth"
	ordermod "github.com/t-line/backend/internal/modules/order"
	"github.com/t-line/backend/internal/modules/payment"
)

// walletAdapter implements payment.WalletOperator using auth.Repository.
type walletAdapter struct {
	repo *authmod.Repository
}

func newWalletAdapter(repo *authmod.Repository) *walletAdapter {
	return &walletAdapter{repo: repo}
}

func (w *walletAdapter) GetBalance(ctx context.Context, userID int64) (decimal.Decimal, error) {
	wallet, err := w.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return decimal.Zero, err
	}
	return wallet.Balance, nil
}

func (w *walletAdapter) DeductBalance(ctx context.Context, userID int64, amount decimal.Decimal, orderID int64, remark string) error {
	wallet, err := w.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet.Balance.LessThan(amount) {
		return fmt.Errorf("余额不足")
	}

	wallet.Balance = wallet.Balance.Sub(amount)
	if err := w.repo.UpdateWalletWithVersion(ctx, wallet); err != nil {
		return fmt.Errorf("扣减余额失败: %w", err)
	}

	_ = w.repo.CreateWalletTransaction(ctx, &authmod.WalletTransaction{
		WalletID:       wallet.ID,
		Type:           "deduct",
		Amount:         amount.Neg(),
		BalanceAfter:   wallet.Balance,
		RelatedOrderID: &orderID,
		Remark:         remark,
		CreatedAt:      time.Now(),
	})

	return nil
}

func (w *walletAdapter) FreezeBalance(ctx context.Context, userID int64, amount decimal.Decimal) error {
	wallet, err := w.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	available := wallet.Balance.Sub(wallet.FrozenAmount)
	if available.LessThan(amount) {
		return fmt.Errorf("可用余额不足")
	}

	wallet.FrozenAmount = wallet.FrozenAmount.Add(amount)
	return w.repo.UpdateWalletWithVersion(ctx, wallet)
}

func (w *walletAdapter) UnfreezeBalance(ctx context.Context, userID int64, amount decimal.Decimal) error {
	wallet, err := w.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	wallet.FrozenAmount = wallet.FrozenAmount.Sub(amount)
	if wallet.FrozenAmount.LessThan(decimal.Zero) {
		wallet.FrozenAmount = decimal.Zero
	}
	return w.repo.UpdateWalletWithVersion(ctx, wallet)
}

func (w *walletAdapter) RefundBalance(ctx context.Context, userID int64, amount decimal.Decimal, orderID int64, remark string) error {
	wallet, err := w.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	wallet.Balance = wallet.Balance.Add(amount)
	if err := w.repo.UpdateWalletWithVersion(ctx, wallet); err != nil {
		return fmt.Errorf("退款失败: %w", err)
	}

	_ = w.repo.CreateWalletTransaction(ctx, &authmod.WalletTransaction{
		WalletID:       wallet.ID,
		Type:           "refund",
		Amount:         amount,
		BalanceAfter:   wallet.Balance,
		RelatedOrderID: &orderID,
		Remark:         remark,
		CreatedAt:      time.Now(),
	})

	return nil
}

// wechatPayAdapter implements payment.WechatPayer using wechat.PayClient.
type wechatPayAdapter struct {
	client wechatPayClient
}

// wechatPayClient is the interface the actual wechat.PayClient satisfies.
type wechatPayClient interface {
	CreatePrepayOrder(orderNo string, amount decimal.Decimal, description string, openID string) (*payment.WechatPayParams, error)
	RefundOrder(paymentNo string, refundNo string, totalAmount, refundAmount decimal.Decimal) error
}

func newWechatPayAdapter(client wechatPayClient) *wechatPayAdapter {
	return &wechatPayAdapter{client: client}
}

func (a *wechatPayAdapter) CreatePrepayOrder(orderNo string, amount decimal.Decimal, description string, openID string) (*payment.WechatPayParams, error) {
	return a.client.CreatePrepayOrder(orderNo, amount, description, openID)
}

func (a *wechatPayAdapter) RefundOrder(paymentNo string, refundNo string, totalAmount, refundAmount decimal.Decimal) error {
	return a.client.RefundOrder(paymentNo, refundNo, totalAmount, refundAmount)
}

// orderQuerierAdapter implements payment.OrderQuerier using order.Service.
type orderQuerierAdapter struct {
	svc *ordermod.Service
}

func newOrderQuerierAdapter(svc *ordermod.Service) *orderQuerierAdapter {
	return &orderQuerierAdapter{svc: svc}
}

func (a *orderQuerierAdapter) GetOrderByID(ctx context.Context, orderID int64) (payment.OrderInfo, error) {
	order, err := a.svc.GetOrderByID(ctx, orderID)
	if err != nil {
		return payment.OrderInfo{}, err
	}
	return payment.OrderInfo{
		ID:        order.ID,
		UserID:    order.UserID,
		Status:    order.Status,
		PayAmount: order.PayAmount,
		Type:      order.Type,
	}, nil
}

// bookingOrderCreator implements booking.OrderCreator using order.Service.
type bookingOrderCreator struct {
	svc *ordermod.Service
}

func newBookingOrderCreator(svc *ordermod.Service) *bookingOrderCreator {
	return &bookingOrderCreator{svc: svc}
}

func (a *bookingOrderCreator) CreateOrderForBooking(ctx context.Context, userID int64, bookingID int64, venueName string, date, startTime, endTime string, amount decimal.Decimal) (int64, error) {
	expiresAt := time.Now().Add(30 * time.Minute)

	itemName := fmt.Sprintf("%s %s %s-%s", venueName, date, startTime, endTime)

	order := &ordermod.Order{
		OrderNo:        ordermod.GenerateOrderNo(),
		UserID:         userID,
		Type:           ordermod.TypeBooking,
		Status:         ordermod.StatusPending,
		TotalAmount:    amount,
		DiscountAmount: decimal.Zero,
		PayAmount:      amount,
		BalancePaid:    decimal.Zero,
		WechatPaid:     decimal.Zero,
		ExpiresAt:      &expiresAt,
		Items: []ordermod.OrderItem{
			{
				ItemType:   "booking",
				ItemID:     bookingID,
				ItemName:   itemName,
				Quantity:   1,
				UnitPrice:  amount,
				TotalPrice: amount,
			},
		},
	}

	if err := a.svc.CreateOrderDirect(ctx, order); err != nil {
		return 0, err
	}
	return order.ID, nil
}
