import { post, get } from './request'

export const paymentApi = {
  preparePayment: (data: {
    order_id: number
    pay_method: 'balance' | 'wechat'
    coupon_id?: number
  }) => post('/payments/prepare', data),

  getCoupons: () =>
    get('/coupons'),

  getAvailableCoupons: (orderType: string, amount: number) =>
    get('/coupons/available', { order_type: orderType, amount }),
}
