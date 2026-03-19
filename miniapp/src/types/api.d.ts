export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T = any> {
  list: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface LoginResp {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

export interface UserInfo {
  id: number
  phone: string
  nickname: string
  avatar_url: string
  gender: number
  age: number
  utr_rating: string | null
  ball_age: number
  self_level: string
  member_level: number
  member_level_name: string
  role: string
}

export interface WalletInfo {
  balance: string
  frozen_amount: string
  total_recharged: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  expires_in: number
}
