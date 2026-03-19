<template>
  <view class="wallet">
    <!-- 余额大卡 -->
    <view class="balance-card">
      <text class="balance-label">账户余额（元）</text>
      <text class="balance-amount">¥{{ balance }}</text>
      <view class="recharge-btn" @tap="showRechargePanel = true">
        <text class="recharge-btn-text">充值</text>
      </view>
    </view>

    <!-- 充值套餐 -->
    <view class="section" v-if="showRechargePanel">
      <text class="section-title">充值套餐</text>
      <view class="package-grid">
        <view
          class="package-item"
          v-for="pkg in packages"
          :key="pkg.amount"
          :class="{ 'package-item--active': selectedPackage === pkg.amount }"
          @tap="selectedPackage = pkg.amount"
        >
          <text class="package-amount">¥{{ pkg.amount }}</text>
          <view class="package-bonus" v-if="pkg.bonus > 0">
            <text>送¥{{ pkg.bonus }}</text>
          </view>
        </view>
      </view>
      <view class="recharge-confirm-btn" @tap="onRecharge">
        <text class="recharge-confirm-text">立即充值 ¥{{ selectedPackage }}</text>
      </view>
    </view>

    <!-- 交易记录 -->
    <view class="section">
      <text class="section-title">交易记录</text>
      <view class="transaction-list">
        <view class="transaction-item" v-for="item in transactions" :key="item.id">
          <view class="transaction-left">
            <text class="transaction-desc">{{ item.description }}</text>
            <text class="transaction-time">{{ item.time }}</text>
          </view>
          <text
            class="transaction-amount"
            :class="{
              'amount-income': item.amount > 0,
              'amount-expense': item.amount < 0,
            }"
          >
            {{ item.amount > 0 ? '+' : '' }}{{ item.amount.toFixed(2) }}
          </text>
        </view>

        <view class="empty" v-if="transactions.length === 0">
          <text class="empty-text">暂无交易记录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'

interface RechargePackage {
  amount: number
  bonus: number
}

interface Transaction {
  id: number
  description: string
  time: string
  amount: number
}

const userStore = useUserStore()

const balance = ref('888.00')
const showRechargePanel = ref(false)
const selectedPackage = ref(300)

const packages = ref<RechargePackage[]>([
  { amount: 100, bonus: 0 },
  { amount: 300, bonus: 20 },
  { amount: 500, bonus: 50 },
  { amount: 1000, bonus: 150 },
])

// ========== Mock 交易记录 ==========
const transactions = ref<Transaction[]>([
  { id: 1, description: '余额充值', time: '2026-03-19 14:30', amount: 300 },
  { id: 2, description: '1号场预约 - 03/18 18:00', time: '2026-03-18 10:15', amount: -280 },
  { id: 3, description: '取消退款 - 4号场', time: '2026-03-15 16:40', amount: 140 },
  { id: 4, description: '3号场预约 - 03/14 10:00', time: '2026-03-14 08:30', amount: -140 },
  { id: 5, description: '余额充值', time: '2026-03-10 20:00', amount: 500 },
  { id: 6, description: '练习场A预约 - 03/08 16:00', time: '2026-03-08 12:00', amount: -80 },
])

// ========== API 调用（保留结构） ==========
async function loadWallet() {
  try {
    // TODO: 接入真实 API
    // await userStore.fetchWallet()
    // balance.value = userStore.walletInfo?.balance || '0.00'
  } catch {
    uni.showToast({ title: '加载钱包失败', icon: 'none' })
  }
}

async function onRecharge() {
  try {
    // TODO: 接入真实 API（微信支付）
    // await paymentApi.preparePayment({
    //   order_id: 0,
    //   pay_method: 'wechat',
    // })
    uni.showToast({ title: '充值功能开发中', icon: 'none' })
  } catch {
    uni.showToast({ title: '充值失败', icon: 'none' })
  }
}

onMounted(() => {
  loadWallet()
})
</script>

<style lang="scss" scoped>
.wallet {
  background: $color-bg;
  min-height: 100vh;
}

// ========== 余额卡 ==========
.balance-card {
  background: linear-gradient(135deg, $brand-primary, $brand-primary-light);
  margin: $spacing-sm $spacing-base;
  border-radius: $radius-lg;
  padding: $spacing-lg $spacing-base;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-sm;
}

.balance-label {
  font-size: $font-size-sm;
  color: rgba(255, 255, 255, 0.8);
}

.balance-amount {
  font-size: 72rpx;
  font-weight: 700;
  color: #fff;
}

.recharge-btn {
  background: rgba(255, 255, 255, 0.25);
  padding: $spacing-xs $spacing-lg;
  border-radius: $radius-round;
  margin-top: $spacing-xs;
}

.recharge-btn-text {
  font-size: $font-size-base;
  color: #fff;
  font-weight: 600;
}

// ========== 通用 Section ==========
.section {
  margin: $spacing-sm $spacing-base;
  background: $color-bg-card;
  border-radius: $radius-lg;
  padding: $spacing-base;
}

.section-title {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-text;
  margin-bottom: $spacing-sm;
}

// ========== 充值套餐 ==========
.package-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-sm;
}

.package-item {
  background: $color-bg;
  border-radius: $radius-base;
  padding: $spacing-base $spacing-sm;
  text-align: center;
  position: relative;
  border: 2rpx solid transparent;
}

.package-item--active {
  border-color: $brand-primary;
  background: #EEF2FF;
}

.package-amount {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-text;
}

.package-bonus {
  position: absolute;
  top: 0;
  right: 0;
  background: $color-error;
  color: #fff;
  font-size: $font-size-xs;
  padding: 2rpx 12rpx;
  border-radius: 0 $radius-base 0 $radius-base;
}

.recharge-confirm-btn {
  margin-top: $spacing-base;
  background: $brand-primary;
  border-radius: $radius-round;
  padding: $spacing-base 0;
  text-align: center;
}

.recharge-confirm-text {
  color: #fff;
  font-size: $font-size-base;
  font-weight: 600;
}

// ========== 交易记录 ==========
.transaction-list {
  display: flex;
  flex-direction: column;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-sm 0;
  border-bottom: 1rpx solid $color-border;

  &:last-child {
    border-bottom: none;
  }
}

.transaction-left {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.transaction-desc {
  font-size: $font-size-base;
  color: $color-text;
}

.transaction-time {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
}

.transaction-amount {
  font-size: $font-size-base;
  font-weight: 700;
}

.amount-income {
  color: $color-success;
}

.amount-expense {
  color: $color-error;
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
