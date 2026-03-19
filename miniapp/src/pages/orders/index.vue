<template>
  <view class="orders">
    <!-- 顶部 Tab -->
    <scroll-view scroll-x class="tabs-scroll">
      <view class="tabs">
        <view
          class="tab"
          v-for="tab in orderTabs"
          :key="tab.value"
          :class="{ 'tab--active': activeTab === tab.value }"
          @tap="activeTab = tab.value"
        >
          <text>{{ tab.label }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 订单列表 -->
    <view class="order-list">
      <view class="order-card" v-for="item in filteredOrders" :key="item.id">
        <view class="card-header">
          <text class="order-no">订单号: {{ item.orderNo }}</text>
          <view
            class="status-tag"
            :class="{
              'status-pending': item.status === 'pending',
              'status-paid': item.status === 'paid',
              'status-used': item.status === 'used',
              'status-completed': item.status === 'completed',
              'status-refunded': item.status === 'refunded',
            }"
          >
            <text>{{ item.statusText }}</text>
          </view>
        </view>

        <view class="card-body">
          <text class="order-name">{{ item.name }}</text>
          <text class="order-detail">{{ item.detail }}</text>
        </view>

        <view class="card-footer">
          <text class="order-price">¥{{ item.price }}</text>
          <view class="order-actions">
            <view
              class="action-btn action-btn--outline"
              v-if="item.status === 'pending'"
              @tap="onCancelOrder(item.id)"
            >
              <text>取消</text>
            </view>
            <view
              class="action-btn action-btn--primary"
              v-if="item.status === 'pending'"
              @tap="onPayOrder(item.id)"
            >
              <text>去支付</text>
            </view>
            <view
              class="action-btn action-btn--primary"
              v-if="item.status === 'completed' || item.status === 'used'"
              @tap="onRebook(item)"
            >
              <text>再次预订</text>
            </view>
          </view>
        </view>
      </view>

      <view class="empty" v-if="filteredOrders.length === 0">
        <text class="empty-text">暂无相关订单</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { orderApi } from '@/api/order'

interface OrderItem {
  id: number
  orderNo: string
  name: string
  detail: string
  price: number
  status: 'pending' | 'paid' | 'used' | 'completed' | 'refunded'
  statusText: string
}

const orderTabs = [
  { label: '全部', value: 'all' },
  { label: '待支付', value: 'pending' },
  { label: '待使用', value: 'paid' },
  { label: '已完成', value: 'completed' },
  { label: '退款', value: 'refunded' },
]

const activeTab = ref('all')

// ========== Mock 数据 ==========
const orders = ref<OrderItem[]>([
  {
    id: 1,
    orderNo: 'TL20260320001',
    name: '1号场预约',
    detail: '2026-03-21 14:00-16:00',
    price: 280,
    status: 'pending',
    statusText: '待支付',
  },
  {
    id: 2,
    orderNo: 'TL20260319002',
    name: '3号场预约',
    detail: '2026-03-22 10:00-11:00',
    price: 140,
    status: 'paid',
    statusText: '待使用',
  },
  {
    id: 3,
    orderNo: 'TL20260318003',
    name: '2号场预约',
    detail: '2026-03-18 18:00-20:00',
    price: 360,
    status: 'completed',
    statusText: '已完成',
  },
  {
    id: 4,
    orderNo: 'TL20260315004',
    name: '4号场预约',
    detail: '2026-03-15 09:00-10:00',
    price: 140,
    status: 'refunded',
    statusText: '已退款',
  },
  {
    id: 5,
    orderNo: 'TL20260312005',
    name: '练习场A预约',
    detail: '2026-03-12 16:00-17:00',
    price: 80,
    status: 'used',
    statusText: '已使用',
  },
])

const filteredOrders = computed(() => {
  if (activeTab.value === 'all') return orders.value
  if (activeTab.value === 'completed') {
    return orders.value.filter(
      (o) => o.status === 'completed' || o.status === 'used',
    )
  }
  return orders.value.filter((o) => o.status === activeTab.value)
})

// ========== 操作 ==========
async function onCancelOrder(id: number) {
  uni.showModal({
    title: '取消订单',
    content: '确定要取消该订单吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        // TODO: 接入真实 API
        // await orderApi.cancelOrder(id)
        orders.value = orders.value.map((o) =>
          o.id === id ? { ...o, status: 'refunded' as const, statusText: '已退款' } : o,
        )
        uni.showToast({ title: '已取消', icon: 'success' })
      } catch {
        uni.showToast({ title: '取消失败', icon: 'none' })
      }
    },
  })
}

function onPayOrder(id: number) {
  // TODO: 跳转支付
  uni.showToast({ title: '支付功能开发中', icon: 'none' })
}

function onRebook(item: OrderItem) {
  uni.navigateTo({ url: '/pages/booking/index' })
}
</script>

<style lang="scss" scoped>
.orders {
  background: $color-bg;
  min-height: 100vh;
}

// ========== Tabs ==========
.tabs-scroll {
  background: $color-bg-card;
  white-space: nowrap;
  border-bottom: 1rpx solid $color-border;
}

.tabs {
  display: inline-flex;
  padding: 0 $spacing-sm;
}

.tab {
  padding: $spacing-base $spacing-sm;
  font-size: $font-size-base;
  color: $color-text-secondary;
  position: relative;
  flex-shrink: 0;
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
    width: 60rpx;
    height: 6rpx;
    border-radius: 3rpx;
    background: $brand-primary;
  }
}

// ========== 订单列表 ==========
.order-list {
  padding: $spacing-sm $spacing-base;
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.order-card {
  background: $color-bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-base $spacing-base $spacing-xs;
}

.order-no {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
}

// ========== 状态标签 ==========
.status-tag {
  font-size: $font-size-xs;
  padding: 2rpx 12rpx;
  border-radius: $radius-sm;
}

.status-pending {
  background: #FEF3C7;
  color: $color-warning;
}

.status-paid {
  background: #EEF2FF;
  color: $brand-primary;
}

.status-used,
.status-completed {
  background: #ECFDF5;
  color: $color-success;
}

.status-refunded {
  background: $color-bg;
  color: $color-text-placeholder;
}

// ========== 卡片内容 ==========
.card-body {
  padding: $spacing-xs $spacing-base;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.order-name {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.order-detail {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

// ========== 卡片底部 ==========
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-sm $spacing-base;
  border-top: 1rpx solid $color-border;
}

.order-price {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-error;
}

.order-actions {
  display: flex;
  gap: $spacing-sm;
}

.action-btn {
  padding: $spacing-xs $spacing-base;
  border-radius: $radius-round;
  font-size: $font-size-sm;
}

.action-btn--outline {
  border: 1rpx solid $color-border;
  color: $color-text-secondary;
}

.action-btn--primary {
  background: $brand-primary;
  color: #fff;
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
