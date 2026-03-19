<template>
  <view class="profile">
    <!-- 用户信息头部 -->
    <view class="header" :style="{ paddingTop: statusBarHeight + 'px' }">
      <view class="user-info" v-if="userStore.isLoggedIn">
        <image class="avatar" :src="userStore.avatarUrl || '/static/images/avatar-default.png'" mode="aspectFill" />
        <view class="user-detail">
          <text class="nickname">{{ userStore.nickname }}</text>
          <view class="badges">
            <view class="member-badge">{{ userStore.memberLevelName }}</view>
            <view class="utr-badge" v-if="userStore.userInfo?.utr_rating">
              UTR {{ userStore.userInfo.utr_rating }}
            </view>
          </view>
        </view>
      </view>
      <view class="user-info" v-else @tap="goLogin">
        <view class="avatar avatar-placeholder">
          <text>👤</text>
        </view>
        <view class="user-detail">
          <text class="nickname">点击登录</text>
          <text class="login-tip">登录享受更多权益</text>
        </view>
      </view>
    </view>

    <!-- 钱包区域 -->
    <view class="wallet-bar" v-if="userStore.isLoggedIn">
      <view class="wallet-item" @tap="navigateTo('/pages/wallet/index')">
        <text class="wallet-value">{{ walletBalance }}</text>
        <text class="wallet-label">余额</text>
      </view>
      <view class="wallet-divider" />
      <view class="wallet-item">
        <text class="wallet-value">{{ pointsBalance }}</text>
        <text class="wallet-label">积分</text>
      </view>
      <view class="wallet-divider" />
      <view class="wallet-item">
        <text class="wallet-value">{{ couponCount }}</text>
        <text class="wallet-label">优惠券</text>
      </view>
    </view>

    <!-- 订单入口 -->
    <view class="order-bar" v-if="userStore.isLoggedIn">
      <view class="order-header">
        <text class="order-title">我的订单</text>
        <text class="order-all" @tap="navigateTo('/pages/orders/index')">全部订单</text>
      </view>
      <view class="order-icons">
        <view class="order-item" v-for="item in orderTypes" :key="item.text" @tap="navigateTo(item.url)">
          <view class="order-icon-wrap">
            <text class="order-icon">{{ item.icon }}</text>
            <view class="order-badge" v-if="item.count > 0">{{ item.count }}</view>
          </view>
          <text class="order-label">{{ item.text }}</text>
        </view>
      </view>
    </view>

    <!-- 功能列表 -->
    <view class="menu-group">
      <view class="menu-item" v-for="item in menuItems" :key="item.text" @tap="navigateTo(item.url)">
        <text class="menu-icon">{{ item.icon }}</text>
        <text class="menu-text">{{ item.text }}</text>
        <text class="menu-arrow">›</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const statusBarHeight = ref(44)

uni.getSystemInfo({ success: (info) => { statusBarHeight.value = info.statusBarHeight || 44 } })

const walletBalance = ref('0.00')
const pointsBalance = ref('0')
const couponCount = ref('0')

const orderTypes = ref([
  { icon: '💳', text: '待支付', count: 2, url: '/pages/orders/index?status=pending' },
  { icon: '📋', text: '待使用', count: 1, url: '/pages/orders/index?status=paid' },
  { icon: '✅', text: '已完成', count: 0, url: '/pages/orders/index?status=completed' },
  { icon: '↩️', text: '退款', count: 0, url: '/pages/orders/index?status=refund' },
])

const menuItems = ref([
  { icon: '📅', text: '我的预约', url: '/pages/my-bookings/index' },
  { icon: '📚', text: '我的课程', url: '/pages/my-courses/index' },
  { icon: '👑', text: '会员中心', url: '/pages/member/index' },
  { icon: '🔔', text: '消息通知', url: '/pages/settings/index' },
  { icon: '⚙️', text: '设置', url: '/pages/settings/index' },
  { icon: '❓', text: '帮助中心', url: '/pages/settings/index' },
])

onMounted(async () => {
  if (userStore.isLoggedIn) {
    await userStore.fetchWallet()
    if (userStore.walletInfo) {
      walletBalance.value = userStore.walletInfo.balance
    }
  }
})

function navigateTo(url: string) {
  if (!userStore.isLoggedIn) {
    goLogin()
    return
  }
  uni.navigateTo({ url })
}

function goLogin() {
  userStore.wechatLogin()
}
</script>

<style lang="scss" scoped>
.profile { padding-bottom: 120rpx; }

.header {
  background: linear-gradient(135deg, #2255CC, #4477DD);
  padding-bottom: 40rpx;
}
.user-info {
  display: flex; align-items: center; gap: 24rpx; padding: 32rpx 32rpx 0;
}
.avatar { width: 120rpx; height: 120rpx; border-radius: 50%; border: 4rpx solid rgba(255,255,255,0.5); }
.avatar-placeholder {
  background: rgba(255,255,255,0.2); display: flex; align-items: center;
  justify-content: center; font-size: 48rpx;
}
.user-detail { flex: 1; }
.nickname { font-size: 36rpx; font-weight: 700; color: #fff; }
.login-tip { font-size: 24rpx; color: rgba(255,255,255,0.7); margin-top: 8rpx; }
.badges { display: flex; gap: 12rpx; margin-top: 12rpx; }
.member-badge {
  background: linear-gradient(135deg, #FFD700, #FFA500); color: #fff;
  font-size: 20rpx; padding: 4rpx 16rpx; border-radius: 999rpx; font-weight: 600;
}
.utr-badge {
  background: rgba(255,255,255,0.2); color: #fff;
  font-size: 20rpx; padding: 4rpx 16rpx; border-radius: 999rpx;
}

.wallet-bar {
  display: flex; align-items: center; justify-content: space-around;
  background: #fff; margin: -20rpx 24rpx 0; padding: 32rpx;
  border-radius: 16rpx; position: relative; z-index: 1;
  box-shadow: 0 4rpx 16rpx rgba(0,0,0,0.06);
}
.wallet-item { display: flex; flex-direction: column; align-items: center; gap: 8rpx; }
.wallet-value { font-size: 36rpx; font-weight: 700; color: #333; }
.wallet-label { font-size: 24rpx; color: #999; }
.wallet-divider { width: 1rpx; height: 60rpx; background: #e5e7eb; }

.order-bar {
  background: #fff; margin: 16rpx 24rpx; border-radius: 16rpx; padding: 24rpx;
}
.order-header { display: flex; justify-content: space-between; margin-bottom: 24rpx; }
.order-title { font-size: 30rpx; font-weight: 700; color: #333; }
.order-all { font-size: 24rpx; color: #2255CC; }
.order-icons { display: flex; justify-content: space-around; }
.order-item { display: flex; flex-direction: column; align-items: center; gap: 8rpx; }
.order-icon-wrap { position: relative; }
.order-icon { font-size: 48rpx; }
.order-badge {
  position: absolute; top: -8rpx; right: -16rpx; min-width: 32rpx; height: 32rpx;
  background: #EF4444; color: #fff; font-size: 20rpx; border-radius: 999rpx;
  display: flex; align-items: center; justify-content: center; padding: 0 8rpx;
}
.order-label { font-size: 24rpx; color: #666; }

.menu-group {
  background: #fff; margin: 16rpx 24rpx; border-radius: 16rpx; overflow: hidden;
}
.menu-item {
  display: flex; align-items: center; padding: 28rpx 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
  &:last-child { border-bottom: none; }
}
.menu-icon { font-size: 36rpx; margin-right: 20rpx; }
.menu-text { flex: 1; font-size: 28rpx; color: #333; }
.menu-arrow { font-size: 32rpx; color: #ccc; }
</style>
