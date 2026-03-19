import { get, post } from './request'

export const orderApi = {
  getOrders: (params?: any) => get('/admin/orders', params),
  reviewRefund: (id: number, data: any) => post(`/admin/refunds/${id}/review`, data),
}
