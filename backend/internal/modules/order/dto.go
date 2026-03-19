package order

import "time"

// --- Request DTOs ---

type CreateOrderReq struct {
	Type   string           `json:"type" binding:"required,oneof=booking product activity recharge"`
	Items  []CreateItemReq  `json:"items" binding:"required,min=1"`
	Remark string           `json:"remark" binding:"omitempty,max=512"`
}

type CreateItemReq struct {
	ItemType  string `json:"item_type" binding:"required"`
	ItemID    int64  `json:"item_id" binding:"required"`
	ItemName  string `json:"item_name" binding:"required,max=128"`
	SkuID     *int64 `json:"sku_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	UnitPrice string `json:"unit_price" binding:"required"`
}

type CancelOrderReq struct {
	Reason string `json:"reason" binding:"omitempty,max=512"`
}

type RefundReviewReq struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Remark string `json:"remark" binding:"omitempty,max=512"`
}

// --- Response DTOs ---

type OrderResp struct {
	ID             int64          `json:"id"`
	OrderNo        string         `json:"order_no"`
	UserID         int64          `json:"user_id"`
	Type           string         `json:"type"`
	Status         string         `json:"status"`
	TotalAmount    string         `json:"total_amount"`
	DiscountAmount string         `json:"discount_amount"`
	PayAmount      string         `json:"pay_amount"`
	BalancePaid    string         `json:"balance_paid"`
	WechatPaid     string         `json:"wechat_paid"`
	CouponID       *int64         `json:"coupon_id"`
	Remark         string         `json:"remark"`
	PaidAt         *time.Time     `json:"paid_at"`
	ExpiresAt      *time.Time     `json:"expires_at"`
	CreatedAt      time.Time      `json:"created_at"`
	Items          []OrderItemResp `json:"items"`
}

type OrderItemResp struct {
	ID         int64  `json:"id"`
	ItemType   string `json:"item_type"`
	ItemID     int64  `json:"item_id"`
	ItemName   string `json:"item_name"`
	Quantity   int    `json:"quantity"`
	UnitPrice  string `json:"unit_price"`
	TotalPrice string `json:"total_price"`
}

type RefundResp struct {
	ID            int64      `json:"id"`
	RefundNo      string     `json:"refund_no"`
	OrderID       int64      `json:"order_id"`
	Amount        string     `json:"amount"`
	BalanceRefund string     `json:"balance_refund"`
	WechatRefund  string     `json:"wechat_refund"`
	Reason        string     `json:"reason"`
	Status        string     `json:"status"`
	ReviewedBy    *int64     `json:"reviewed_by"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func ToOrderResp(o *Order) OrderResp {
	items := make([]OrderItemResp, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, OrderItemResp{
			ID:         item.ID,
			ItemType:   item.ItemType,
			ItemID:     item.ItemID,
			ItemName:   item.ItemName,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice.StringFixed(2),
			TotalPrice: item.TotalPrice.StringFixed(2),
		})
	}

	return OrderResp{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		UserID:         o.UserID,
		Type:           o.Type,
		Status:         o.Status,
		TotalAmount:    o.TotalAmount.StringFixed(2),
		DiscountAmount: o.DiscountAmount.StringFixed(2),
		PayAmount:      o.PayAmount.StringFixed(2),
		BalancePaid:    o.BalancePaid.StringFixed(2),
		WechatPaid:     o.WechatPaid.StringFixed(2),
		CouponID:       o.CouponID,
		Remark:         o.Remark,
		PaidAt:         o.PaidAt,
		ExpiresAt:      o.ExpiresAt,
		CreatedAt:      o.CreatedAt,
		Items:          items,
	}
}

func ToRefundResp(r *Refund) RefundResp {
	return RefundResp{
		ID:            r.ID,
		RefundNo:      r.RefundNo,
		OrderID:       r.OrderID,
		Amount:        r.Amount.StringFixed(2),
		BalanceRefund: r.BalanceRefund.StringFixed(2),
		WechatRefund:  r.WechatRefund.StringFixed(2),
		Reason:        r.Reason,
		Status:        r.Status,
		ReviewedBy:    r.ReviewedBy,
		ReviewedAt:    r.ReviewedAt,
		CreatedAt:     r.CreatedAt,
	}
}
