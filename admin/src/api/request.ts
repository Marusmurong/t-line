import axios from 'axios'
import type { AxiosRequestConfig, AxiosResponse } from 'axios'

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

const instance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

// request interceptor
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// response interceptor
instance.interceptors.response.use(
  (res: AxiosResponse<ApiResponse>) => {
    const data = res.data
    if (data.code === 0) {
      return data.data
    }
    return Promise.reject(new Error(data.message || '请求失败'))
  },
  async (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

export function get<T = any>(url: string, params?: any) {
  return instance.get<any, T>(url, { params })
}

export function post<T = any>(url: string, data?: any) {
  return instance.post<any, T>(url, data)
}

export function put<T = any>(url: string, data?: any) {
  return instance.put<any, T>(url, data)
}

export function del<T = any>(url: string, params?: any) {
  return instance.delete<any, T>(url, { params })
}

export default instance
