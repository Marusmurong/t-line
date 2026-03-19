import type { ApiResponse } from '@/types/api'

const BASE_URL = 'http://localhost:8080/api/v1'

let isRefreshing = false
let pendingRequests: Array<() => void> = []

function getToken(): string {
  return uni.getStorageSync('access_token') || ''
}

function getRefreshToken(): string {
  return uni.getStorageSync('refresh_token') || ''
}

function setTokens(access: string, refresh: string) {
  uni.setStorageSync('access_token', access)
  uni.setStorageSync('refresh_token', refresh)
}

function clearTokens() {
  uni.removeStorageSync('access_token')
  uni.removeStorageSync('refresh_token')
}

async function refreshToken(): Promise<boolean> {
  try {
    const res = await uni.request({
      url: `${BASE_URL}/auth/refresh`,
      method: 'POST',
      data: { refresh_token: getRefreshToken() },
    })
    const data = (res.data as ApiResponse).data
    if (data?.access_token) {
      setTokens(data.access_token, data.refresh_token)
      return true
    }
    return false
  } catch {
    return false
  }
}

export function request<T = any>(options: {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  data?: any
  header?: Record<string, string>
  auth?: boolean
}): Promise<T> {
  const { url, method = 'GET', data, header = {}, auth = true } = options

  return new Promise((resolve, reject) => {
    const token = getToken()
    if (auth && token) {
      header['Authorization'] = `Bearer ${token}`
    }

    uni.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        ...header,
      },
      success: async (res) => {
        const statusCode = res.statusCode
        const body = res.data as ApiResponse<T>

        if (statusCode === 200 && body.code === 0) {
          resolve(body.data)
          return
        }

        // token expired, try refresh
        if (statusCode === 401) {
          if (!isRefreshing) {
            isRefreshing = true
            const ok = await refreshToken()
            isRefreshing = false

            if (ok) {
              pendingRequests.forEach((cb) => cb())
              pendingRequests = []
              // retry original request
              try {
                const result = await request<T>(options)
                resolve(result)
              } catch (e) {
                reject(e)
              }
              return
            }

            // refresh failed, redirect to login
            clearTokens()
            uni.showToast({ title: '登录已过期', icon: 'none' })
            reject(new Error('token expired'))
            return
          }

          // queue request while refreshing
          return new Promise<void>((r) => {
            pendingRequests.push(() => r())
          }).then(() => request<T>(options).then(resolve).catch(reject))
        }

        // business error
        uni.showToast({ title: body.message || '请求失败', icon: 'none' })
        reject(new Error(body.message))
      },
      fail: (err) => {
        uni.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      },
    })
  })
}

export const get = <T = any>(url: string, data?: any) =>
  request<T>({ url, method: 'GET', data })

export const post = <T = any>(url: string, data?: any) =>
  request<T>({ url, method: 'POST', data })

export const put = <T = any>(url: string, data?: any) =>
  request<T>({ url, method: 'PUT', data })

export const del = <T = any>(url: string, data?: any) =>
  request<T>({ url, method: 'DELETE', data })
