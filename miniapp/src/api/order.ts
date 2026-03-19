import { get, post } from './request'

export const orderApi = {
  getOrders: (params?: { status?: string; page?: number }) =>
    get('/orders', params),

  getOrderDetail: (id: number) =>
    get(`/orders/${id}`),

  cancelOrder: (id: number) =>
    post(`/orders/${id}/cancel`),
}
