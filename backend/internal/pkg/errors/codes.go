package errors

import "fmt"

// AppError is a structured error with error code.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// General (40000-40099)
var (
	ErrInvalidParams = New(40000, "参数错误")
	ErrRecordNotFound = New(40001, "记录不存在")
	ErrDuplicateEntry = New(40002, "记录已存在")
)

// Auth (40100-40199)
var (
	ErrUnauthorized     = New(40100, "未授权")
	ErrTokenExpired     = New(40101, "登录已过期")
	ErrTokenInvalid     = New(40102, "无效令牌")
	ErrForbidden        = New(40103, "权限不足")
	ErrUserDisabled     = New(40104, "账户已禁用")
	ErrSMSCodeInvalid   = New(40105, "验证码错误或已过期")
	ErrSMSTooFrequent   = New(40106, "验证码发送过于频繁")
	ErrPhoneRegistered  = New(40107, "手机号已注册")
	ErrPasswordWrong    = New(40108, "密码错误")
	ErrWeChatLoginFail  = New(40109, "微信登录失败")
)

// Booking (40200-40299)
var (
	ErrSlotUnavailable  = New(40200, "该时段不可预约")
	ErrSlotConflict     = New(40201, "时段冲突")
	ErrBookingNotFound  = New(40202, "预约不存在")
	ErrBookingCancelled = New(40203, "预约已取消")
	ErrCancelTimeout    = New(40204, "已超过取消时限")
	ErrWaitlistFull     = New(40205, "候补队列已满")
)

// Payment (40300-40399)
var (
	ErrInsufficientBalance = New(40300, "余额不足")
	ErrPaymentFailed       = New(40301, "支付失败")
	ErrPaymentExpired      = New(40302, "支付已超时")
	ErrRefundFailed        = New(40303, "退款失败")
	ErrCouponInvalid       = New(40304, "优惠券无效")
	ErrCouponExpired       = New(40305, "优惠券已过期")
	ErrOrderPaid           = New(40306, "订单已支付")
)

// Product (40400-40499)
var (
	ErrProductNotFound  = New(40400, "商品不存在")
	ErrProductOffShelf  = New(40401, "商品已下架")
	ErrStockInsufficient = New(40402, "库存不足")
)

// Activity (40500-40599)
var (
	ErrActivityNotFound   = New(40500, "活动不存在")
	ErrActivityFull       = New(40501, "活动已满员")
	ErrActivityClosed     = New(40502, "活动报名已截止")
	ErrAlreadyRegistered  = New(40503, "已报名该活动")
)

// Academic (40600-40699)
var (
	ErrScheduleConflict    = New(40600, "排课时间冲突")
	ErrCoachUnavailable    = New(40601, "教练不可用")
	ErrVenueUnavailable    = New(40602, "场地不可用")
	ErrCourseNotFound      = New(40603, "课程不存在")
)

// Server (50000-50099)
var (
	ErrInternal        = New(50000, "服务器内部错误")
	ErrDatabaseFail    = New(50001, "数据库错误")
	ErrCacheFail       = New(50002, "缓存错误")
	ErrExternalService = New(50003, "外部服务错误")
)
