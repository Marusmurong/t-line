<template>
  <view class="home">
    <!-- 搜索栏 -->
    <view class="header" :style="{ paddingTop: statusBarHeight + 'px' }">
      <view class="header-inner">
        <image class="logo" src="/static/images/logo.png" mode="aspectFit" />
        <view class="search-bar" @tap="onSearch">
          <text class="search-icon">🔍</text>
          <text class="search-placeholder">搜索场地、活动、课程</text>
        </view>
      </view>
    </view>

    <!-- Banner -->
    <swiper class="banner" indicator-dots autoplay circular :interval="4000">
      <swiper-item v-for="item in banners" :key="item.id">
        <view class="banner-item" :style="{ background: item.bg }">
          <text class="banner-title">{{ item.title }}</text>
          <text class="banner-desc">{{ item.desc }}</text>
        </view>
      </swiper-item>
    </swiper>

    <!-- 快捷入口 -->
    <view class="shortcuts">
      <view class="shortcut-item" v-for="item in shortcuts" :key="item.text" @tap="navigateTo(item.url)">
        <view class="shortcut-icon" :style="{ background: item.color }">
          <text>{{ item.icon }}</text>
        </view>
        <text class="shortcut-text">{{ item.text }}</text>
      </view>
    </view>

    <!-- 今日场地 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">今日场地</text>
        <text class="section-more" @tap="navigateTo('/pages/booking/index')">查看全部</text>
      </view>
      <view class="venue-summary">
        <view class="venue-stat">
          <text class="venue-stat-num">4</text>
          <text class="venue-stat-label">可预约</text>
        </view>
        <view class="venue-divider" />
        <view class="venue-stat">
          <text class="venue-stat-num">6</text>
          <text class="venue-stat-label">总场地</text>
        </view>
        <view class="venue-divider" />
        <view class="venue-stat">
          <text class="venue-stat-num venue-stat-highlight">67%</text>
          <text class="venue-stat-label">使用率</text>
        </view>
      </view>
    </view>

    <!-- 热门活动 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">热门活动</text>
        <text class="section-more" @tap="navigateTo('/pages/activity/index')">更多</text>
      </view>
      <scroll-view scroll-x class="activity-scroll">
        <view class="activity-card" v-for="item in hotActivities" :key="item.id" @tap="navigateTo('/pages/activity/detail?id=' + item.id)">
          <view class="activity-cover" :style="{ background: item.bg }">
            <view class="activity-tag">{{ item.type }}</view>
          </view>
          <text class="activity-name">{{ item.title }}</text>
          <text class="activity-time">{{ item.time }}</text>
          <text class="activity-price">¥{{ item.price }}</text>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const statusBarHeight = ref(0)

// get status bar height
uni.getSystemInfo({
  success: (info) => {
    statusBarHeight.value = info.statusBarHeight || 44
  },
})

const banners = ref([
  { id: 1, title: '春季特惠', desc: '充值 500 送 100', bg: 'linear-gradient(135deg, #2255CC, #4477DD)' },
  { id: 2, title: '周末团课', desc: '每周六日精彩团课', bg: 'linear-gradient(135deg, #22C55E, #4ADE80)' },
  { id: 3, title: '新会员福利', desc: '注册即享首次订场8折', bg: 'linear-gradient(135deg, #F59E0B, #FBBF24)' },
])

const shortcuts = ref([
  { icon: '🎾', text: '订场', url: '/pages/booking/index', color: '#EEF2FF' },
  { icon: '🏆', text: '活动', url: '/pages/activity/index', color: '#FEF3C7' },
  { icon: '📚', text: '课程', url: '/pages/shop/index', color: '#ECFDF5' },
  { icon: '🛒', text: '商城', url: '/pages/shop/index', color: '#FEF2F2' },
])

const hotActivities = ref([
  { id: 1, title: '周六畅打', type: '畅打', time: '03/22 14:00', price: '68', bg: '#2255CC' },
  { id: 2, title: '初级团课', type: '团课', time: '03/23 10:00', price: '128', bg: '#22C55E' },
  { id: 3, title: '春季联赛', type: '赛事', time: '04/01 09:00', price: '200', bg: '#F59E0B' },
])

function navigateTo(url: string) {
  uni.navigateTo({ url })
}

function onSearch() {
  uni.showToast({ title: '搜索功能开发中', icon: 'none' })
}
</script>

<style lang="scss" scoped>
.home { padding-bottom: 120rpx; }

.header {
  background: #fff;
  .header-inner {
    display: flex; align-items: center; padding: 16rpx 24rpx; gap: 20rpx;
  }
}
.logo { width: 80rpx; height: 80rpx; }
.search-bar {
  flex: 1; display: flex; align-items: center; gap: 12rpx;
  background: #f5f6f8; border-radius: 999rpx; padding: 16rpx 24rpx;
}
.search-placeholder { color: #999; font-size: 26rpx; }

.banner {
  height: 300rpx; margin: 16rpx 24rpx; border-radius: 16rpx; overflow: hidden;
}
.banner-item {
  height: 100%; display: flex; flex-direction: column;
  justify-content: center; padding: 40rpx; border-radius: 16rpx;
}
.banner-title { color: #fff; font-size: 36rpx; font-weight: 700; }
.banner-desc { color: rgba(255,255,255,0.85); font-size: 26rpx; margin-top: 12rpx; }

.shortcuts {
  display: flex; justify-content: space-around;
  background: #fff; margin: 16rpx 24rpx; padding: 32rpx 0;
  border-radius: 16rpx;
}
.shortcut-item { display: flex; flex-direction: column; align-items: center; gap: 12rpx; }
.shortcut-icon {
  width: 96rpx; height: 96rpx; border-radius: 24rpx;
  display: flex; align-items: center; justify-content: center; font-size: 40rpx;
}
.shortcut-text { font-size: 24rpx; color: #333; }

.section { margin: 16rpx 24rpx; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16rpx; }
.section-title { font-size: 32rpx; font-weight: 700; color: #333; }
.section-more { font-size: 24rpx; color: #2255CC; }

.venue-summary {
  display: flex; align-items: center; justify-content: space-around;
  background: #fff; border-radius: 16rpx; padding: 32rpx;
}
.venue-stat { display: flex; flex-direction: column; align-items: center; gap: 8rpx; }
.venue-stat-num { font-size: 44rpx; font-weight: 700; color: #2255CC; }
.venue-stat-highlight { color: #C8E632; }
.venue-stat-label { font-size: 24rpx; color: #999; }
.venue-divider { width: 1rpx; height: 60rpx; background: #e5e7eb; }

.activity-scroll { white-space: nowrap; }
.activity-card {
  display: inline-flex; flex-direction: column;
  width: 280rpx; margin-right: 16rpx;
  background: #fff; border-radius: 16rpx; overflow: hidden;
}
.activity-cover {
  height: 160rpx; display: flex; align-items: flex-start; justify-content: flex-end;
  padding: 12rpx;
}
.activity-tag {
  background: rgba(255,255,255,0.9); color: #2255CC; font-size: 20rpx;
  padding: 4rpx 12rpx; border-radius: 8rpx; font-weight: 600;
}
.activity-name { font-size: 28rpx; font-weight: 600; padding: 12rpx 16rpx 4rpx; }
.activity-time { font-size: 22rpx; color: #999; padding: 0 16rpx; }
.activity-price { font-size: 28rpx; color: #EF4444; font-weight: 700; padding: 8rpx 16rpx 16rpx; }
</style>
