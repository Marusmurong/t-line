import { get } from './request'

export const productApi = {
  getProducts: (params?: any) =>
    get('/products', params),

  getProductDetail: (id: number) =>
    get(`/products/${id}`),
}
