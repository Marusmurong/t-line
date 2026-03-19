<template>
  <view class="confirm">
    <!-- 场地信息卡 -->
    <view class="venue-card">
      <view class="venue-card-bg">
        <text class="venue-card-name">{{ bookingData.venueName }}</text>
        <view class="venue-card-tag">
          <text>{{ venueTypeLabel(bookingData.venueType) }}</text>
        </view>
      </view>
    </view>

    <!-- 预订时段 -->
    <view class="section-card">
      <text class="section-title">预订时段</text>
      <view class="slot-list">
        <view class="slot-item" v-for="(slot, idx) in bookingData.slots" :key="idx">
          <view class="slot-left">
            <text class="slot-date">{{ bookingData.date }}</text>
            <text class="slot-time">{{ slot.time }} - {{ nextHour(slot.time) }}</text>
          </view>
          <view class="slot-right">
            <text class="slot-duration">1小时</text>
            <text class="slot-price">¥{{ slot.price }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 费用明细 -->
    <view class="section-card">
      <text class="section-title">费用明细</text>
      <view class="fee-row">
        <text class="fee-label">场地费</text>
        <text class="fee-value">¥{{ bookingData.totalPrice }}</text>
      </view>
      <view class="fee-row">
        <text class="fee-label">会员折扣</text>
        <text class="fee-value fee-discount">-¥{{ discount }}</text>
      </view>
      <view class="fee-divider" />
      <view class="fee-row fee-total">
        <text class="fee-label">实付金额</text>
        <text class="fee-value fee-total-value">¥{{ finalPrice }}</text>
      </view>
    </view>

    <!-- 支付方式 -->
    <view class="section-card">
      <text class="section-title">支付方式</text>
      <view class="pay-option" @tap="onSelectPayMethod('balance')">
        <view class="pay-option-left">
          <text class="pay-icon">💰</text>
          <view class="pay-info">
            <text class="pay-name">余额支付</text>
            <text class="pay-desc">余额 ¥{{ walletBalance }}</text>
          </view>
        </view>
        <view class="radio" :class="{ 'radio--checked': payMethod === 'balance' }" />
      </view>
      <view class="pay-option" @tap="onSelectPayMethod('wechat')">
        <view class="pay-option-left">
          <text class="pay-icon">💳</text>
          <view class="pay-info">
            <text class="pay-name">微信支付</text>
          </view>
        </view>
        <view class="radio" :class="{ 'radio--checked': payMethod === 'wechat' }" />
      </view>
    </view>

    <!-- 订场须知 -->
    <view class="section-card">
      <view class="notice-header" @tap="noticeExpanded = !noticeExpanded">
        <text class="section-title">订场须知</text>
        <text class="notice-arrow">{{ noticeExpanded ? '收起' : '展开' }}</text>
      </view>
      <view class="notice-content" v-if="noticeExpanded">
        <text class="notice-text">1. 预订成功后，请在预约时段开始前15分钟到达场地签到。</text>
        <text class="notice-text">2. 如需取消预订，请在开场前2小时操作，逾期将扣除50%费用。</text>
        <text class="notice-text">3. 请穿着专业网球鞋入场，禁止穿着黑底鞋。</text>
        <text class="notice-text">4. 场地提供免费饮用水，球拍租赁需另付费。</text>
        <text class="notice-text">5. 请爱护场地设施，损坏需按价赔偿。</text>
      </view>
    </view>

    <!-- 底部按钮 -->
    <view class="bottom-bar">
      <view class="confirm-btn" @tap="onConfirmPay">
        <text class="confirm-btn-text">确认支付 ¥{{ finalPrice }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onLoad } from 'vue'
import { bookingApi } from '@/api/booking'
import { paymentApi } from '@/api/payment'
import { useUserStore } from '@/stores/user'

interface SlotInfo {
  venueId: number
  venueName: string
  venueType: string
  time: string
  price: number
}

interface BookingData {
  date: string
  venueName: string
  venueType: string
  slots: SlotInfo[]
  totalPrice: number
}

const userStore = useUserStore()

const bookingData = ref<BookingData>({
  date: '',
  venueName: '',
  venueType: '',
  slots: [],
  totalPrice: 0,
})

const payMethod = ref<'balance' | 'wechat'>('balance')
const noticeExpanded = ref(false)
const walletBalance = ref('888.00')
const isSubmitting = ref(false)

// 模拟会员折扣 10%
const discount = computed(() => Math.round(bookingData.value.totalPrice * 0.1))
const finalPrice = computed(() => bookingData.value.totalPrice - discount.value)

onLoad((options: Record<string, string> | undefined) => {
  if (options?.data) {
    try {
      const parsed = JSON.parse(decodeURIComponent(options.data))
      const firstSlot = parsed.slots[0] || {}
      bookingData.value = {
        date: parsed.date,
        venueName: firstSlot.venueName || '',
        venueType: firstSlot.venueType || '',
        slots: parsed.slots,
        totalPrice: parsed.totalPrice,
      }
    } catch {
      uni.showToast({ title: '参数解析失败', icon: 'none' })
    }
  }

  // 获取钱包余额
  if (userStore.walletInfo) {
    walletBalance.value = userStore.walletInfo.balance
  }
})

function venueTypeLabel(type: string): string {
  const map: Record<string, string> = {
    indoor: '室内',
    outdoor: '室外',
    practice: '练习',
  }
  return map[type] || type
}

function nextHour(time: string): string {
  const hour = parseInt(time.split(':')[0], 10)
  return `${String(hour + 1).padStart(2, '0')}:00`
}

function onSelectPayMethod(method: 'balance' | 'wechat') {
  payMethod.value = method
}

async function onConfirmPay() {
  if (isSubmitting.value) return
  isSubmitting.value = true

  try {
    // TODO: 接入真实 API
    // const booking = await bookingApi.createBooking({
    //   venue_id: bookingData.value.slots[0]?.venueId,
    //   date: bookingData.value.date,
    //   time_slots: bookingData.value.slots.map(s => s.time),
    // })
    // await paymentApi.preparePayment({
    //   order_id: booking.id,
    //   pay_method: payMethod.value,
    // })

    uni.showToast({ title: '预订成功', icon: 'success' })
    setTimeout(() => {
      uni.redirectTo({ url: '/pages/my-bookings/index' })
    }, 1500)
  } catch {
    uni.showToast({ title: '支付失败，请重试', icon: 'none' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.confirm {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 160rpx;
}

// ========== 场地卡片 ==========
.venue-card {
  margin: $spacing-sm $spacing-base;
}

.venue-card-bg {
  background: linear-gradient(135deg, $brand-primary, $brand-primary-light);
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-base;
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.venue-card-name {
  font-size: $font-size-xl;
  font-weight: 700;
  color: #fff;
}

.venue-card-tag {
  background: rgba(255, 255, 255, 0.25);
  padding: 4rpx 16rpx;
  border-radius: $radius-sm;
  font-size: $font-size-sm;
  color: #fff;
}

// ========== 通用卡片 ==========
.section-card {
  background: $color-bg-card;
  margin: $spacing-sm $spacing-base;
  border-radius: $radius-lg;
  padding: $spacing-base;
}

.section-title {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-text;
  margin-bottom: $spacing-sm;
}

// ========== 时段列表 ==========
.slot-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.slot-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-sm;
  background: $color-bg;
  border-radius: $radius-base;
}

.slot-left {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.slot-date {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.slot-time {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.slot-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4rpx;
}

.slot-duration {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
}

.slot-price {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

// ========== 费用明细 ==========
.fee-row {
  display: flex;
  justify-content: space-between;
  padding: $spacing-xs 0;
}

.fee-label {
  font-size: $font-size-base;
  color: $color-text-secondary;
}

.fee-value {
  font-size: $font-size-base;
  color: $color-text;
}

.fee-discount {
  color: $color-success;
}

.fee-divider {
  height: 1rpx;
  background: $color-border;
  margin: $spacing-sm 0;
}

.fee-total {
  .fee-label {
    font-weight: 700;
    color: $color-text;
  }
}

.fee-total-value {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-error;
}

// ========== 支付方式 ==========
.pay-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-sm 0;
  border-bottom: 1rpx solid $color-border;

  &:last-child {
    border-bottom: none;
  }
}

.pay-option-left {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.pay-icon {
  font-size: $font-size-xxl;
}

.pay-info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.pay-name {
  font-size: $font-size-base;
  color: $color-text;
}

.pay-desc {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
}

.radio {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 3rpx solid $color-border;
}

.radio--checked {
  border-color: $brand-primary;
  background: $brand-primary;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 16rpx;
    height: 16rpx;
    border-radius: 50%;
    background: #fff;
  }
}

// ========== 订场须知 ==========
.notice-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notice-arrow {
  font-size: $font-size-sm;
  color: $brand-primary;
}

.notice-content {
  display: flex;
  flex-direction: column;
  gap: $spacing-xs;
  margin-top: $spacing-sm;
}

.notice-text {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: 1.6;
}

// ========== 底部按钮 ==========
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: $spacing-sm $spacing-base;
  padding-bottom: calc(#{$spacing-sm} + env(safe-area-inset-bottom));
  background: $color-bg-card;
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.06);
}

.confirm-btn {
  background: $brand-primary;
  border-radius: $radius-round;
  padding: $spacing-base 0;
  text-align: center;
}

.confirm-btn-text {
  color: #fff;
  font-size: $font-size-md;
  font-weight: 700;
}
</style>
