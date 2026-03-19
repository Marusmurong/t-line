import { post, get, put } from './request'
import type { LoginResp, UserInfo, WalletInfo, TokenPair } from '@/types/api'

export const authApi = {
  wechatLogin: (code: string) =>
    post<LoginResp>('/auth/wechat-login', { code }),

  phoneLogin: (phone: string, code: string) =>
    post<LoginResp>('/auth/phone-login', { phone, code }),

  passwordLogin: (phone: string, password: string) =>
    post<LoginResp>('/auth/password-login', { phone, password }),

  sendSMSCode: (phone: string) =>
    post('/auth/sms-code', { phone }),

  refreshToken: (refreshToken: string) =>
    post<TokenPair>('/auth/refresh', { refresh_token: refreshToken }),

  getProfile: () =>
    get<UserInfo>('/auth/profile'),

  updateProfile: (data: Partial<UserInfo>) =>
    put<UserInfo>('/auth/profile', data),

  getWallet: () =>
    get<WalletInfo>('/wallet'),
}
