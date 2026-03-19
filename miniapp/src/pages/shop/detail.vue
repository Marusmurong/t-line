<template>
  <view class="detail">
    <!-- 商品图片区域 -->
    <view class="product-hero">
      <view class="product-img" />
    </view>

    <!-- 价格区域 -->
    <view class="price-section">
      <view class="price-row">
        <text class="price-current">¥{{ product.price }}</text>
        <text class="price-original" v-if="product.originalPrice">¥{{ product.originalPrice }}</text>
      </view>
      <view class="sold-info">
        <text class="sold-text">已售 {{ product.sold }}</text>
      </view>
    </view>

    <!-- 商品信息 -->
    <view class="info-section">
      <text class="product-name">{{ product.name }}</text>
      <text class="product-desc">{{ product.desc }}</text>
    </view>

    <!-- SKU 选择 -->
    <view class="sku-section" v-if="product.skus.length > 0">
      <text class="sku-title">规格选择</text>
      <view class="sku-options">
        <view
          class="sku-option"
          v-for="(sku, idx) in product.skus"
          :key="idx"
          :class="{ 'sku-option--active': selectedSku === idx }"
          @tap="onSelectSku(idx)"
        >
          <text class="sku-option-text">{{ sku }}</text>
        </view>
      </view>
    </view>

    <!-- 数量选择器 -->
    <view class="quantity-section">
      <text class="quantity-label">购买数量</text>
      <view class="quantity-picker">
        <view class="qty-btn" @tap="onChangeQty(-1)">
          <text class="qty-btn-text">-</text>
        </view>
        <text class="qty-value">{{ quantity }}</text>
        <view class="qty-btn" @tap="onChangeQty(1)">
          <text class="qty-btn-text">+</text>
        </view>
      </view>
    </view>

    <!-- 商品详情描述 -->
    <view class="desc-section">
      <text class="desc-title">商品详情</text>
      <text class="desc-content">{{ product.detail }}</text>
    </view>

    <!-- 底部固定操作栏 -->
    <view class="bottom-bar">
      <view class="btn-cart" @tap="onAddCart">
        <text class="btn-cart-text">加入购物车</text>
      </view>
      <view class="btn-buy" @tap="onBuyNow">
        <text class="btn-buy-text">立即购买</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// ========== 类型 ==========
interface Product {
  id: number
  name: string
  desc: string
  detail: string
  price: number
  originalPrice: number
  sold: number
  skus: string[]
}

// ========== 页面参数 ==========
const productId = ref(0)

// ========== Mock 数据 ==========
const product = ref<Product>({
  id: 301,
  name: 'Wilson Pro Staff 97 V14 网球拍',
  desc: '费德勒签名款，碳纤维材质，精准控球利器',
  detail: '产品参数：\n- 拍面：97平方英寸\n- 重量：315g（未穿线）\n- 平衡点：310mm\n- 硬度：63\n- 材质：高模碳纤维\n- 穿线范围：16x19\n\n适合中高级选手，提供出色的控球感和手感反馈。经典Pro Staff家族的最新力作，传承经典设计的同时融入现代科技。',
  price: 1899,
  originalPrice: 2399,
  sold: 128,
  skus: ['4 1/4 (L2)', '4 3/8 (L3)', '4 1/2 (L4)'],
})

const selectedSku = ref(0)
const quantity = ref(1)

// ========== 接收页面参数 ==========
function onLoad(options: Record<string, string>) {
  if (options.id) {
    productId.value = Number(options.id)
    // TODO: 接入真实 API
    // const data = await productApi.getProductDetail(productId.value)
  }
}

// #ifdef MP-WEIXIN
defineExpose({ onLoad })
// #endif

// ========== SKU 选择 ==========
function onSelectSku(idx: number) {
  selectedSku.value = idx
}

// ========== 数量选择 ==========
function onChangeQty(delta: number) {
  const next = quantity.value + delta
  if (next < 1) return
  if (next > 99) return
  quantity.value = next
}

// ========== 购物车 & 购买 ==========
function onAddCart() {
  uni.showToast({
    title: '已加入购物车',
    icon: 'success',
  })
}

function onBuyNow() {
  uni.showToast({
    title: '购买功能开发中',
    icon: 'none',
  })
}
</script>

<style lang="scss" scoped>
.detail {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 160rpx;
}

// ========== 商品图片 ==========
.product-hero {
  padding: 0;
}

.product-img {
  width: 100%;
  height: 560rpx;
  background: linear-gradient(135deg, $brand-primary, $brand-primary-light);
}

// ========== 价格区域 ==========
.price-section {
  background: $color-bg-card;
  padding: $spacing-base;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: $spacing-sm;
}

.price-current {
  font-size: $font-size-xxl;
  font-weight: 700;
  color: $color-error;
}

.price-original {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
  text-decoration: line-through;
}

.sold-text {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

// ========== 商品信息 ==========
.info-section {
  background: $color-bg-card;
  padding: 0 $spacing-base $spacing-base;
  display: flex;
  flex-direction: column;
  gap: $spacing-xs;
}

.product-name {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-text;
}

.product-desc {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: 1.6;
}

// ========== SKU 选择 ==========
.sku-section {
  background: $color-bg-card;
  margin-top: $spacing-sm;
  padding: $spacing-base;
}

.sku-title {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-sm;
}

.sku-options {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
}

.sku-option {
  padding: $spacing-xs $spacing-base;
  border: 2rpx solid $color-border;
  border-radius: $radius-base;
  transition: all 0.2s;
}

.sku-option--active {
  border-color: $brand-primary;
  background: #EEF2FF;
}

.sku-option-text {
  font-size: $font-size-sm;
  color: $color-text;
}

.sku-option--active .sku-option-text {
  color: $brand-primary;
  font-weight: 600;
}

// ========== 数量选择 ==========
.quantity-section {
  background: $color-bg-card;
  margin-top: $spacing-sm;
  padding: $spacing-base;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.quantity-label {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.quantity-picker {
  display: flex;
  align-items: center;
  gap: 0;
}

.qty-btn {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $color-bg;
  border-radius: $radius-sm;
}

.qty-btn-text {
  font-size: $font-size-lg;
  color: $color-text;
  font-weight: 600;
}

.qty-value {
  width: 80rpx;
  text-align: center;
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

// ========== 商品详情 ==========
.desc-section {
  background: $color-bg-card;
  margin-top: $spacing-sm;
  padding: $spacing-base;
}

.desc-title {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-sm;
}

.desc-content {
  font-size: $font-size-sm;
  color: $color-text-secondary;
  line-height: 1.8;
  white-space: pre-wrap;
}

// ========== 底部操作栏 ==========
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: $spacing-sm;
  padding: $spacing-sm $spacing-base;
  padding-bottom: calc(#{$spacing-sm} + env(safe-area-inset-bottom));
  background: $color-bg-card;
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.06);
}

.btn-cart {
  flex: 1;
  background: $color-bg;
  border: 2rpx solid $brand-primary;
  padding: $spacing-sm 0;
  border-radius: $radius-round;
  text-align: center;
}

.btn-cart-text {
  color: $brand-primary;
  font-size: $font-size-base;
  font-weight: 600;
}

.btn-buy {
  flex: 1;
  background: $brand-primary;
  padding: $spacing-sm 0;
  border-radius: $radius-round;
  text-align: center;
}

.btn-buy-text {
  color: #fff;
  font-size: $font-size-base;
  font-weight: 600;
}
</style>
