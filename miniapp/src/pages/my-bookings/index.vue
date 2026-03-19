<template>
  <view class="my-bookings">
    <!-- 顶部 Tab -->
    <view class="tabs">
      <view
        class="tab"
        :class="{ 'tab--active': activeTab === 'upcoming' }"
        @tap="activeTab = 'upcoming'"
      >
        <text>即将开始</text>
      </view>
      <view
        class="tab"
        :class="{ 'tab--active': activeTab === 'history' }"
        @tap="activeTab = 'history'"
      >
        <text>历史记录</text>
      </view>
    </view>

    <!-- 即将开始 -->
    <view class="booking-list" v-if="activeTab === 'upcoming'">
      <view class="booking-card" v-for="item in upcomingList" :key="item.id">
        <view class="card-header">
          <view class="card-venue">
            <text class="card-venue-name">{{ item.venueName }}</text>
            <view class="status-tag status-paid">
              <text>{{ item.statusText }}</text>
            </view>
          </view>
          <text class="card-countdown">{{ item.countdown }}</text>
        </view>
        <view class="card-body">
          <view class="card-info-row">
            <text class="card-label">日期</text>
            <text class="card-value">{{ item.date }}</text>
          </view>
          <view class="card-info-row">
            <text class="card-label">时间</text>
            <text class="card-value">{{ item.timeRange }}</text>
          </view>
          <view class="card-info-row">
            <text class="card-label">金额</text>
            <text class="card-value card-price">¥{{ item.price }}</text>
          </view>
        </view>
        <view class="card-footer">
          <view class="cancel-btn" @tap="onCancel(item.id)">
            <text class="cancel-btn-text">取消预约</text>
          </view>
        </view>
      </view>

      <view class="empty" v-if="upcomingList.length === 0">
        <text class="empty-text">暂无即将开始的预约</text>
      </view>
    </view>

    <!-- 历史记录 -->
    <view class="booking-list" v-if="activeTab === 'history'">
      <view
        class="booking-card booking-card--history"
        v-for="item in historyList"
        :key="item.id"
      >
        <view class="card-header">
          <view class="card-venue">
            <text class="card-venue-name">{{ item.venueName }}</text>
            <view
              class="status-tag"
              :class="{
                'status-completed': item.status === 'completed',
                'status-cancelled': item.status === 'cancelled',
              }"
            >
              <text>{{ item.statusText }}</text>
            </view>
          </view>
        </view>
        <view class="card-body">
          <view class="card-info-row">
            <text class="card-label">日期</text>
            <text class="card-value">{{ item.date }}</text>
          </view>
          <view class="card-info-row">
            <text class="card-label">时间</text>
            <text class="card-value">{{ item.timeRange }}</text>
          </view>
          <view class="card-info-row">
            <text class="card-label">金额</text>
            <text class="card-value">¥{{ item.price }}</text>
          </view>
        </view>
      </view>

      <view class="empty" v-if="historyList.length === 0">
        <text class="empty-text">暂无历史记录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { bookingApi } from '@/api/booking'

interface BookingItem {
  id: number
  venueName: string
  date: string
  timeRange: string
  price: number
  status: 'paid' | 'completed' | 'cancelled'
  statusText: string
  countdown: string
}

const activeTab = ref<'upcoming' | 'history'>('upcoming')

// ========== Mock 数据 ==========
const upcomingList = ref<BookingItem[]>([
  {
    id: 1,
    venueName: '1号场（室内）',
    date: '2026-03-21',
    timeRange: '14:00 - 16:00',
    price: 280,
    status: 'paid',
    statusText: '已支付',
    countdown: '距开场 2小时30分',
  },
  {
    id: 2,
    venueName: '3号场（室外）',
    date: '2026-03-22',
    timeRange: '10:00 - 11:00',
    price: 140,
    status: 'paid',
    statusText: '已支付',
    countdown: '距开场 1天22小时',
  },
])

const historyList = ref<BookingItem[]>([
  {
    id: 3,
    venueName: '2号场（室内）',
    date: '2026-03-18',
    timeRange: '18:00 - 20:00',
    price: 360,
    status: 'completed',
    statusText: '已完成',
    countdown: '',
  },
  {
    id: 4,
    venueName: '4号场（室外）',
    date: '2026-03-15',
    timeRange: '09:00 - 10:00',
    price: 140,
    status: 'cancelled',
    statusText: '已取消',
    countdown: '',
  },
  {
    id: 5,
    venueName: '练习场A',
    date: '2026-03-12',
    timeRange: '16:00 - 17:00',
    price: 80,
    status: 'completed',
    statusText: '已完成',
    countdown: '',
  },
])

// ========== API 调用（保留结构） ==========
async function loadBookings() {
  try {
    // TODO: 接入真实 API
    // const upcoming = await bookingApi.getBookings({ status: 'upcoming' })
    // const history = await bookingApi.getBookings({ status: 'history' })
  } catch {
    uni.showToast({ title: '加载预约失败', icon: 'none' })
  }
}

async function onCancel(id: number) {
  uni.showModal({
    title: '取消预约',
    content: '确定要取消该预约吗？开场前2小时内取消将扣除50%费用。',
    success: async (res) => {
      if (!res.confirm) return
      try {
        // TODO: 接入真实 API
        // await bookingApi.cancelBooking(id)
        upcomingList.value = upcomingList.value.filter((item) => item.id !== id)
        uni.showToast({ title: '已取消', icon: 'success' })
      } catch {
        uni.showToast({ title: '取消失败', icon: 'none' })
      }
    },
  })
}

onMounted(() => {
  loadBookings()
})
</script>

<style lang="scss" scoped>
.my-bookings {
  background: $color-bg;
  min-height: 100vh;
}

// ========== Tabs ==========
.tabs {
  display: flex;
  background: $color-bg-card;
  border-bottom: 1rpx solid $color-border;
}

.tab {
  flex: 1;
  text-align: center;
  padding: $spacing-base 0;
  font-size: $font-size-base;
  color: $color-text-secondary;
  position: relative;
}

.tab--active {
  color: $brand-primary;
  font-weight: 700;

  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 80rpx;
    height: 6rpx;
    border-radius: 3rpx;
    background: $brand-primary;
  }
}

// ========== 预约列表 ==========
.booking-list {
  padding: $spacing-sm $spacing-base;
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.booking-card {
  background: $color-bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
}

.booking-card--history {
  opacity: 0.75;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-base $spacing-base $spacing-xs;
}

.card-venue {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.card-venue-name {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-text;
}

.card-countdown {
  font-size: $font-size-sm;
  color: $brand-primary;
  font-weight: 600;
}

// ========== 状态标签 ==========
.status-tag {
  font-size: $font-size-xs;
  padding: 2rpx 12rpx;
  border-radius: $radius-sm;
}

.status-paid {
  background: #ECFDF5;
  color: $color-success;
}

.status-completed {
  background: $color-bg;
  color: $color-text-placeholder;
}

.status-cancelled {
  background: #FEF2F2;
  color: $color-error;
}

// ========== 卡片内容 ==========
.card-body {
  padding: $spacing-xs $spacing-base;
  display: flex;
  flex-direction: column;
  gap: $spacing-xs;
}

.card-info-row {
  display: flex;
  justify-content: space-between;
}

.card-label {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
}

.card-value {
  font-size: $font-size-sm;
  color: $color-text;
}

.card-price {
  font-weight: 600;
  color: $color-error;
}

// ========== 卡片底部 ==========
.card-footer {
  display: flex;
  justify-content: flex-end;
  padding: $spacing-sm $spacing-base;
  border-top: 1rpx solid $color-border;
}

.cancel-btn {
  border: 1rpx solid $color-border;
  padding: $spacing-xs $spacing-base;
  border-radius: $radius-round;
}

.cancel-btn-text {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

// ========== 空状态 ==========
.empty {
  display: flex;
  justify-content: center;
  padding: $spacing-xl 0;
}

.empty-text {
  font-size: $font-size-base;
  color: $color-text-placeholder;
}
</style>
