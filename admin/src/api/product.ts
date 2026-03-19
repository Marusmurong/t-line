import { get, post, put, del } from './request'

export const productApi = {
  getProducts: (params?: any) => get('/admin/products', params),
  createProduct: (data: any) => post('/admin/products', data),
  updateProduct: (id: number, data: any) => put(`/admin/products/${id}`, data),
  deleteProduct: (id: number) => del(`/admin/products/${id}`),
}
