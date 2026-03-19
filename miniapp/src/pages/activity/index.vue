<template>
  <view class="activity-page">
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

    <!-- 活动卡片列表 -->
    <scroll-view scroll-y class="activity-list">
      <view
        class="activity-card"
        v-for="item in filteredActivities"
        :key="item.id"
        @tap="goDetail(item.id)"
      >
        <!-- 左侧封面 -->
        <view class="card-cover" :style="{ background: item.coverBg }">
          <view class="card-type-tag">
            <text class="card-type-text">{{ item.type }}</text>
          </view>
        </view>

        <!-- 右侧信息 -->
        <view class="card-body">
          <text class="card-title">{{ item.title }}</text>

          <view class="card-meta">
            <text class="meta-text">{{ item.time }}</text>
          </view>
          <view class="card-meta">
            <text class="meta-text">{{ item.venue }}</text>
          </view>

          <!-- 水平要求标签 -->
          <view class="level-tag" :class="levelClass(item.level)">
            <text class="level-text">{{ item.level }}</text>
          </view>

          <!-- 报名进度 -->
          <view class="progress-row">
            <view class="progress-bar">
              <view
                class="progress-fill"
                :style="{ width: progressPercent(item) + '%' }"
                :class="{ 'progress-full': item.enrolled >= item.capacity }"
              />
            </view>
            <text class="progress-text">{{ item.enrolled }}/{{ item.capacity }}人</text>
          </view>

          <!-- 价格 + 按钮 -->
          <view class="card-footer">
            <text class="card-price">¥{{ item.price }}</text>
            <view
              class="enroll-btn"
              :class="{ 'enroll-btn--disabled': item.enrolled >= item.capacity }"
              @tap.stop="onEnroll(item)"
            >
              <text class="enroll-btn-text">
                {{ item.enrolled >= item.capacity ? '已满' : '报名' }}
              </text>
            </view>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="empty" v-if="filteredActivities.length === 0">
        <text class="empty-text">暂无相关活动</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// ========== 类型 ==========
interface Activity {
  id: number
  title: string
  type: '畅打' | '团课' | '赛事' | '主题活动'
  time: string
  venue: string
  level: '不限' | '初级' | '中级' | '高级'
  enrolled: number
  capacity: number
  price: number
  coverBg: string
}

// ========== Tab ==========
const tabs = [
  { label: '全部', value: 'all' },
  { label: '畅打', value: '畅打' },
  { label: '团课', value: '团课' },
  { label: '赛事', value: '赛事' },
  { label: '主题活动', value: '主题活动' },
]
const activeTab = ref('all')

function onTabChange(value: string) {
  activeTab.value = value
}

// ========== Mock 数据 ==========
const activities = ref<Activity[]>([
  {
    id: 1, title: '周六下午畅打', type: '畅打',
    time: '03/22 14:00-17:00', venue: '1号室内场',
    level: '不限', enrolled: 12, capacity: 16, price: 68,
    coverBg: 'linear-gradient(135deg, #2255CC, #4477DD)',
  },
  {
    id: 2, title: '初级网球团课', type: '团课',
    time: '03/23 10:00-11:30', venue: '3号室外场',
    level: '初级', enrolled: 8, capacity: 8, price: 128,
    coverBg: 'linear-gradient(135deg, #22C55E, #4ADE80)',
  },
  {
    id: 3, title: '春季单打联赛 第一轮', type: '赛事',
    time: '04/01 09:00-18:00', venue: '全场地',
    level: '中级', enrolled: 24, capacity: 32, price: 200,
    coverBg: 'linear-gradient(135deg, #F59E0B, #FBBF24)',
  },
  {
    id: 4, title: '中级提高班团课', type: '团课',
    time: '03/24 14:00-15:30', venue: '2号室内场',
    level: '中级', enrolled: 6, capacity: 10, price: 158,
    coverBg: 'linear-gradient(135deg, #8B5CF6, #A78BFA)',
  },
  {
    id: 5, title: '亲子网球体验日', type: '主题活动',
    time: '03/29 09:00-12:00', venue: '3号/4号室外场',
    level: '不限', enrolled: 18, capacity: 20, price: 99,
    coverBg: 'linear-gradient(135deg, #EC4899, #F472B6)',
  },
  {
    id: 6, title: '高级实战对抗训练', type: '畅打',
    time: '03/30 16:00-18:00', venue: '1号室内场',
    level: '高级', enrolled: 4, capacity: 8, price: 88,
    coverBg: 'linear-gradient(135deg, #EF4444, #F87171)',
  },
  {
    id: 7, title: '企业团建网球活动', type: '主题活动',
    time: '04/05 13:00-17:00', venue: '全场地',
    level: '不限', enrolled: 30, capacity: 40, price: 150,
    coverBg: 'linear-gradient(135deg, #14B8A6, #2DD4BF)',
  },
])

// ========== 筛选 ==========
const filteredActivities = computed(() => {
  if (activeTab.value === 'all') return activities.value
  return activities.value.filter((a) => a.type === activeTab.value)
})

// ========== 进度 ==========
function progressPercent(item: Activity): number {
  return Math.min((item.enrolled / item.capacity) * 100, 100)
}

// ========== 水平要求样式 ==========
function levelClass(level: string): string {
  const map: Record<string, string> = {
    '不限': 'level--any',
    '初级': 'level--beginner',
    '中级': 'level--intermediate',
    '高级': 'level--advanced',
  }
  return map[level] || 'level--any'
}

// ========== 导航 ==========
function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/activity/detail?id=${id}` })
}

function onEnroll(item: Activity) {
  if (item.enrolled >= item.capacity) return
  uni.navigateTo({ url: `/pages/activity/detail?id=${item.id}` })
}
</script>

<style lang="scss" scoped>
.activity-page {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 120rpx;
}

// ========== Tab 栏 ==========
.tab-bar {
  display: flex;
  background: $color-bg-card;
  padding: $spacing-xs $spacing-sm;
  gap: 0;
  position: sticky;
  top: 0;
  z-index: 10;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: $spacing-sm 0;
  font-size: $font-size-sm;
  color: $color-text-secondary;
  border-bottom: 4rpx solid transparent;
  transition: all 0.2s;
}

.tab-item--active {
  color: $brand-primary;
  font-weight: 700;
  border-bottom-color: $brand-primary;
}

// ========== 活动列表 ==========
.activity-list {
  padding: $spacing-sm $spacing-base;
}

// ========== 活动卡片 ==========
.activity-card {
  display: flex;
  background: $color-bg-card;
  border-radius: $radius-lg;
  overflow: hidden;
  margin-bottom: $spacing-sm;
}

.card-cover {
  width: 200rpx;
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding: $spacing-xs;
  position: relative;
}

.card-type-tag {
  background: rgba(255, 255, 255, 0.9);
  padding: 4rpx 12rpx;
  border-radius: $radius-sm;
}

.card-type-text {
  font-size: $font-size-xs;
  font-weight: 600;
  color: $brand-primary;
}

.card-body {
  flex: 1;
  padding: $spacing-sm;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.card-title {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-text;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 6rpx;
}

.meta-text {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

// ========== 水平要求标签 ==========
.level-tag {
  display: inline-flex;
  padding: 2rpx 12rpx;
  border-radius: $radius-sm;
  align-self: flex-start;
}

.level-text {
  font-size: $font-size-xs;
  font-weight: 600;
}

.level--any {
  background: $color-bg;
  .level-text { color: $color-text-secondary; }
}

.level--beginner {
  background: #ECFDF5;
  .level-text { color: $color-success; }
}

.level--intermediate {
  background: #EEF2FF;
  .level-text { color: $brand-primary; }
}

.level--advanced {
  background: #FEF2F2;
  .level-text { color: $color-error; }
}

// ========== 进度条 ==========
.progress-row {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  margin-top: 2rpx;
}

.progress-bar {
  flex: 1;
  height: 12rpx;
  background: $color-bg;
  border-radius: $radius-round;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: $brand-primary;
  border-radius: $radius-round;
  transition: width 0.3s;
}

.progress-full {
  background: $color-text-disabled;
}

.progress-text {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
  flex-shrink: 0;
}

// ========== 底部价格+按钮 ==========
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 4rpx;
}

.card-price {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-error;
}

.enroll-btn {
  background: $brand-primary;
  padding: $spacing-xs $spacing-base;
  border-radius: $radius-round;
}

.enroll-btn-text {
  color: #fff;
  font-size: $font-size-xs;
  font-weight: 600;
}

.enroll-btn--disabled {
  background: $color-text-disabled;
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
