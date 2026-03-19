package payment

import "time"

// --- Request DTOs ---

type PreparePayReq struct {
	OrderID      int64  `json:"order_id" binding:"required"`
	Method       string `json:"method" binding:"required,oneof=balance wechat combo"`
	CouponID     *int64 `json:"coupon_id"`
	BalanceAmount string `json:"balance_amount"` // only for combo pay
}

type WechatCallbackReq struct {
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	TradeState    string `json:"trade_state"`
	Amount        int64  `json:"amount"` // in cents
}

// --- Response DTOs ---

type PreparePayResp struct {
	PaymentNo     string  `json:"payment_no"`
	Method        string  `json:"method"`
	TotalAmount   string  `json:"total_amount"`
	BalanceAmount string  `json:"balance_amount"`
	WechatAmount  string  `json:"wechat_amount"`
	WechatParams  *WechatPayParams `json:"wechat_params,omitempty"`
	Paid          bool    `json:"paid"`
}

type WechatPayParams struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type CouponResp struct {
	ID              int64     `json:"id"`
	UserCouponID    int64     `json:"user_coupon_id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	Value           string    `json:"value"`
	MinAmount       string    `json:"min_amount"`
	ApplicableTypes []string  `json:"applicable_types"`
	Status          int       `json:"status"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type PaymentResp struct {
	ID            int64      `json:"id"`
	PaymentNo     string     `json:"payment_no"`
	OrderID       int64      `json:"order_id"`
	Method        string     `json:"method"`
	Amount        string     `json:"amount"`
	BalanceAmount string     `json:"balance_amount"`
	WechatAmount  string     `json:"wechat_amount"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
}

func ToPaymentResp(p *Payment) PaymentResp {
	return PaymentResp{
		ID:            p.ID,
		PaymentNo:     p.PaymentNo,
		OrderID:       p.OrderID,
		Method:        p.Method,
		Amount:        p.Amount.StringFixed(2),
		BalanceAmount: p.BalanceAmount.StringFixed(2),
		WechatAmount:  p.WechatAmount.StringFixed(2),
		Status:        p.Status,
		PaidAt:        p.PaidAt,
	}
}
