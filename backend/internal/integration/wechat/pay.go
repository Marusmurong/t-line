package wechat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/t-line/backend/internal/modules/payment"
)

// PayClient handles WeChat Pay operations including prepay, callback verification, and refunds.
type PayClient struct {
	mchID     string
	mchAPIKey string
	notifyURL string
}

// NewPayClient creates a new WeChat Pay client.
func NewPayClient(mchID, mchAPIKey, notifyURL string) *PayClient {
	return &PayClient{
		mchID:     mchID,
		mchAPIKey: mchAPIKey,
		notifyURL: notifyURL,
	}
}

// VerifyCallback verifies the callback signature from WeChat Pay.
func (c *PayClient) VerifyCallback(body []byte, signature string) error {
	mac := hmac.New(sha256.New, []byte(c.mchAPIKey))
	mac.Write(body)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSig), []byte(signature)) {
		return fmt.Errorf("wechat pay callback signature mismatch")
	}
	return nil
}

// CreatePrepayOrder creates a prepay order with WeChat Pay and returns JSAPI payment params.
func (c *PayClient) CreatePrepayOrder(orderNo string, amount decimal.Decimal, description string, openID string) (*payment.WechatPayParams, error) {
	// TODO: integrate with actual WeChat Pay API (JSAPI unified order)
	// For now, return placeholder params — the real implementation would:
	// 1. Call POST https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi
	// 2. Get prepay_id from response
	// 3. Build and sign the JSAPI payment parameters

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := uuid.New().String()[:32]

	return &payment.WechatPayParams{
		AppID:     "placeholder",
		TimeStamp: timestamp,
		NonceStr:  nonceStr,
		Package:   "prepay_id=placeholder",
		SignType:  "RSA",
		PaySign:   "placeholder",
	}, nil
}

// RefundOrder initiates a refund through WeChat Pay.
func (c *PayClient) RefundOrder(paymentNo string, refundNo string, totalAmount, refundAmount decimal.Decimal) error {
	// TODO: integrate with actual WeChat Pay refund API
	// POST https://api.mch.weixin.qq.com/v3/refund/domestic/refunds
	return nil
}
