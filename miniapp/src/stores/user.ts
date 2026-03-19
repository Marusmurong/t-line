import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, WalletInfo } from '@/types/api'
import { authApi } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)
  const walletInfo = ref<WalletInfo | null>(null)

  const isLoggedIn = computed(() => !!uni.getStorageSync('access_token'))
  const nickname = computed(() => userInfo.value?.nickname || '未登录')
  const avatarUrl = computed(() => userInfo.value?.avatar_url || '')
  const memberLevelName = computed(() => userInfo.value?.member_level_name || '普通会员')

  async function autoLogin() {
    const token = uni.getStorageSync('access_token')
    if (!token) return

    try {
      const profile = await authApi.getProfile()
      userInfo.value = profile
    } catch {
      // token invalid, clear
      uni.removeStorageSync('access_token')
      uni.removeStorageSync('refresh_token')
    }
  }

  async function wechatLogin() {
    try {
      const [loginErr, loginRes] = await uni.login({ provider: 'weixin' })
      if (loginErr || !loginRes) {
        uni.showToast({ title: '微信登录失败', icon: 'none' })
        return false
      }

      const resp = await authApi.wechatLogin(loginRes.code)
      uni.setStorageSync('access_token', resp.access_token)
      uni.setStorageSync('refresh_token', resp.refresh_token)
      userInfo.value = resp.user
      return true
    } catch {
      return false
    }
  }

  async function phoneLogin(phone: string, code: string) {
    const resp = await authApi.phoneLogin(phone, code)
    uni.setStorageSync('access_token', resp.access_token)
    uni.setStorageSync('refresh_token', resp.refresh_token)
    userInfo.value = resp.user
    return true
  }

  async function fetchProfile() {
    try {
      userInfo.value = await authApi.getProfile()
    } catch {
      // ignore
    }
  }

  async function fetchWallet() {
    try {
      walletInfo.value = await authApi.getWallet()
    } catch {
      // ignore
    }
  }

  function logout() {
    uni.removeStorageSync('access_token')
    uni.removeStorageSync('refresh_token')
    userInfo.value = null
    walletInfo.value = null
    uni.reLaunch({ url: '/pages/home/index' })
  }

  return {
    userInfo,
    walletInfo,
    isLoggedIn,
    nickname,
    avatarUrl,
    memberLevelName,
    autoLogin,
    wechatLogin,
    phoneLogin,
    fetchProfile,
    fetchWallet,
    logout,
  }
})
