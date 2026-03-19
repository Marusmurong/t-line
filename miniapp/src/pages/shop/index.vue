<template>
  <view class="shop">
    <!-- 顶部 Tab -->
    <view class="tab-bar">
      <view
        class="tab-item"
        v-for="tab in tabs"
        :key="tab.value"
        :class="{ 'tab-item--active': activeTab === tab.value }"
        @tap="onTabChange(tab.value)"
      >
        <text>{{ tab.label }}</text>
      </view>
    </view>

    <!-- 课程 Tab -->
    <view v-if="activeTab === 'course'" class="tab-content">
      <!-- 私教课 -->
      <view class="section-label">
        <text class="section-label-text">私教课</text>
      </view>
      <view
        class="coach-card"
        v-for="item in privateCoaches"
        :key="item.id"
        @tap="goDetail(item.id)"
      >
        <view class="coach-avatar" :style="{ background: item.avatarBg }">
          <text class="coach-avatar-text">{{ item.name.charAt(0) }}</text>
        </view>
        <view class="coach-info">
          <text class="coach-name">{{ item.name }}</text>
          <view class="coach-rating">
            <text
              class="star"
              v-for="s in 5"
              :key="s"
              :class="{ 'star--filled': s <= item.rating }"
            >★</text>
            <text class="rating-text">{{ item.rating.toFixed(1) }}</text>
          </view>
          <text class="coach-desc">{{ item.desc }}</text>
        </view>
        <view class="coach-price-area">
          <text class="coach-price">¥{{ item.price }}</text>
          <text class="coach-price-unit">/次</text>
        </view>
      </view>

      <!-- 课包 -->
      <view class="section-label">
        <text class="section-label-text">课包</text>
      </view>
      <view
        class="package-card"
        v-for="item in coursePackages"
        :key="item.id"
        @tap="goDetail(item.id)"
      >
        <view class="package-badge" :style="{ background: item.badgeBg }">
          <text class="package-badge-text">{{ item.badge }}</text>
        </view>
        <text class="package-name">{{ item.name }}</text>
        <view class="package-price-row">
          <text class="package-price">¥{{ item.price }}</text>
          <text class="package-original">¥{{ item.originalPrice }}</text>
        </view>
        <text class="package-save">省 ¥{{ item.originalPrice - item.price }}</text>
      </view>
    </view>

    <!-- 球具 Tab -->
    <view v-if="activeTab === 'equipment'" class="tab-content">
      <view class="product-grid">
        <view
          class="product-item"
          v-for="item in equipments"
          :key="item.id"
          @tap="goDetail(item.id)"
        >
          <view class="product-img" :style="{ background: item.imgBg }">
            <text class="product-img-text">{{ item.name.charAt(0) }}</text>
          </view>
          <view class="product-body">
            <text class="product-name">{{ item.name }}</text>
            <view class="product-bottom">
              <text class="product-price">¥{{ item.price }}</text>
              <view
                class="stock-tag"
                :class="{
                  'stock-tag--available': item.stock === 'available',
                  'stock-tag--preorder': item.stock === 'preorder',
                  'stock-tag--out': item.stock === 'out',
                }"
              >
                <text class="stock-tag-text">{{ stockLabel(item.stock) }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- 服务 Tab -->
    <view v-if="activeTab === 'service'" class="tab-content">
      <view
        class="service-item"
        v-for="item in services"
        :key="item.id"
      >
        <view class="service-icon" :style="{ background: item.iconBg }">
          <text class="service-icon-text">{{ item.icon }}</text>
        </view>
        <view class="service-info">
          <text class="service-name">{{ item.name }}</text>
          <text class="service-desc">{{ item.desc }}</text>
        </view>
        <view class="service-right">
          <text class="service-price">¥{{ item.price }}</text>
          <view class="service-btn" @tap="onBookService(item)">
            <text class="service-btn-text">预约</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// ========== 类型 ==========
interface Coach {
  id: number
  name: string
  avatarBg: string
  rating: number
  desc: string
  price: number
}

interface CoursePackage {
  id: number
  name: string
  badge: string
  badgeBg: string
  price: number
  originalPrice: number
}

interface Equipment {
  id: number
  name: string
  imgBg: string
  price: number
  stock: 'available' | 'preorder' | 'out'
}

interface Service {
  id: number
  name: string
  icon: string
  iconBg: string
  desc: string
  price: number
}

// ========== Tab 切换 ==========
const tabs = [
  { label: '课程', value: 'course' },
  { label: '球具', value: 'equipment' },
  { label: '服务', value: 'service' },
]
const activeTab = ref('course')

function onTabChange(value: string) {
  activeTab.value = value
}

// ========== Mock 数据：私教课 ==========
const privateCoaches = ref<Coach[]>([
  { id: 101, name: '王教练', avatarBg: '#2255CC', rating: 4.8, desc: '前省队选手，10年教学经验', price: 300 },
  { id: 102, name: '李教练', avatarBg: '#22C55E', rating: 4.6, desc: 'PTR认证教练，擅长青少年培训', price: 260 },
  { id: 103, name: '张教练', avatarBg: '#F59E0B', rating: 4.9, desc: 'ITF认证，专攻发球技术提升', price: 350 },
])

// ========== Mock 数据：课包 ==========
const coursePackages = ref<CoursePackage[]>([
  { id: 201, name: '10次私教课包', badge: '热卖', badgeBg: '#EF4444', price: 2680, originalPrice: 3000 },
  { id: 202, name: '20次私教课包', badge: '超值', badgeBg: '#F59E0B', price: 4980, originalPrice: 6000 },
  { id: 203, name: '5次体验课包', badge: '新人', badgeBg: '#22C55E', price: 1280, originalPrice: 1500 },
])

// ========== Mock 数据：球具 ==========
const equipments = ref<Equipment[]>([
  { id: 301, name: 'Wilson Pro Staff', imgBg: 'linear-gradient(135deg, #2255CC, #4477DD)', price: 1899, stock: 'available' },
  { id: 302, name: 'Babolat Pure Aero', imgBg: 'linear-gradient(135deg, #F59E0B, #FBBF24)', price: 1699, stock: 'available' },
  { id: 303, name: 'Head Speed MP', imgBg: 'linear-gradient(135deg, #22C55E, #4ADE80)', price: 1599, stock: 'preorder' },
  { id: 304, name: 'Yonex EZONE 98', imgBg: 'linear-gradient(135deg, #3B82F6, #60A5FA)', price: 1799, stock: 'out' },
  { id: 305, name: 'Luxilon 网球线', imgBg: 'linear-gradient(135deg, #8B5CF6, #A78BFA)', price: 128, stock: 'available' },
  { id: 306, name: 'Wilson 吸汗带', imgBg: 'linear-gradient(135deg, #EC4899, #F472B6)', price: 45, stock: 'available' },
])

// ========== Mock 数据：服务 ==========
const services = ref<Service[]>([
  { id: 401, name: '穿线服务', icon: '🔧', iconBg: '#EEF2FF', desc: '专业穿线师，多种线材可选', price: 40 },
  { id: 402, name: '球拍租赁', icon: '🎾', iconBg: '#ECFDF5', desc: '品牌球拍租赁，按小时计费', price: 30 },
  { id: 403, name: '发球机租赁', icon: '⚡', iconBg: '#FEF3C7', desc: '自动发球机，可调速度角度', price: 60 },
  { id: 404, name: '录像分析', icon: '📹', iconBg: '#FEF2F2', desc: '专业教练录像回放+技术分析', price: 150 },
])

// ========== 库存标签 ==========
function stockLabel(stock: string): string {
  const map: Record<string, string> = {
    available: '有货',
    preorder: '预定',
    out: '缺货',
  }
  return map[stock] || stock
}

// ========== 导航 ==========
function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/shop/detail?id=${id}` })
}

function onBookService(item: Service) {
  uni.showToast({ title: `已预约${item.name}`, icon: 'none' })
}
</script>

<style lang="scss" scoped>
.shop {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 120rpx;
}

// ========== Tab 栏 ==========
.tab-bar {
  display: flex;
  background: $color-bg-card;
  padding: $spacing-xs $spacing-base;
  gap: $spacing-sm;
  position: sticky;
  top: 0;
  z-index: 10;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: $spacing-sm 0;
  font-size: $font-size-base;
  color: $color-text-secondary;
  border-bottom: 4rpx solid transparent;
  transition: all 0.2s;
}

.tab-item--active {
  color: $brand-primary;
  font-weight: 700;
  border-bottom-color: $brand-primary;
}

// ========== Tab 内容区 ==========
.tab-content {
  padding: $spacing-sm $spacing-base;
}

// ========== Section 标签 ==========
.section-label {
  margin: $spacing-sm 0 $spacing-xs;
}

.section-label-text {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-text;
}

// ========== 私教课卡片 ==========
.coach-card {
  display: flex;
  align-items: center;
  background: $color-bg-card;
  border-radius: $radius-lg;
  padding: $spacing-base;
  margin-bottom: $spacing-sm;
  gap: $spacing-base;
}

.coach-avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.coach-avatar-text {
  color: #fff;
  font-size: $font-size-lg;
  font-weight: 700;
}

.coach-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.coach-name {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.coach-rating {
  display: flex;
  align-items: center;
  gap: 2rpx;
}

.star {
  font-size: $font-size-sm;
  color: $color-text-disabled;
}

.star--filled {
  color: $color-warning;
}

.rating-text {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  margin-left: 8rpx;
}

.coach-desc {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

.coach-price-area {
  display: flex;
  align-items: baseline;
  flex-shrink: 0;
}

.coach-price {
  font-size: $font-size-xl;
  font-weight: 700;
  color: $color-error;
}

.coach-price-unit {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

// ========== 课包卡片 ==========
.package-card {
  background: $color-bg-card;
  border-radius: $radius-lg;
  padding: $spacing-base;
  margin-bottom: $spacing-sm;
  position: relative;
  overflow: hidden;
}

.package-badge {
  position: absolute;
  top: 0;
  right: 0;
  padding: 4rpx 20rpx 4rpx 24rpx;
  border-radius: 0 $radius-lg 0 $radius-lg;
}

.package-badge-text {
  color: #fff;
  font-size: $font-size-xs;
  font-weight: 600;
}

.package-name {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-xs;
}

.package-price-row {
  display: flex;
  align-items: baseline;
  gap: $spacing-sm;
}

.package-price {
  font-size: $font-size-xl;
  font-weight: 700;
  color: $color-error;
}

.package-original {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
  text-decoration: line-through;
}

.package-save {
  font-size: $font-size-xs;
  color: $color-success;
  margin-top: 4rpx;
}

// ========== 球具网格 ==========
.product-grid {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
}

.product-item {
  width: calc(50% - #{$spacing-sm} / 2);
  background: $color-bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
}

.product-img {
  height: 240rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-img-text {
  color: rgba(255, 255, 255, 0.6);
  font-size: $font-size-xxl;
  font-weight: 700;
}

.product-body {
  padding: $spacing-sm;
}

.product-name {
  font-size: $font-size-sm;
  font-weight: 600;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-xs;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.product-price {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-error;
}

.stock-tag {
  padding: 2rpx 12rpx;
  border-radius: $radius-sm;
}

.stock-tag-text {
  font-size: $font-size-xs;
}

.stock-tag--available {
  background: #ECFDF5;
  .stock-tag-text { color: $color-success; }
}

.stock-tag--preorder {
  background: #FEF3C7;
  .stock-tag-text { color: $color-warning; }
}

.stock-tag--out {
  background: $color-bg;
  .stock-tag-text { color: $color-text-disabled; }
}

// ========== 服务列表 ==========
.service-item {
  display: flex;
  align-items: center;
  background: $color-bg-card;
  border-radius: $radius-lg;
  padding: $spacing-base;
  margin-bottom: $spacing-sm;
  gap: $spacing-base;
}

.service-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: $radius-base;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.service-icon-text {
  font-size: 40rpx;
}

.service-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.service-name {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.service-desc {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

.service-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: $spacing-xs;
  flex-shrink: 0;
}

.service-price {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-error;
}

.service-btn {
  background: $brand-primary;
  padding: $spacing-xs $spacing-base;
  border-radius: $radius-round;
}

.service-btn-text {
  color: #fff;
  font-size: $font-size-xs;
  font-weight: 600;
}
</style>
