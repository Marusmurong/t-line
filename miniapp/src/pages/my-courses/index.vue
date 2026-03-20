<template>
  <view class="my-courses">
    <!-- 顶部 Tab -->
    <view class="tabs">
      <view
        class="tab"
        :class="{ 'tab--active': activeTab === 'active' }"
        @tap="activeTab = 'active'"
      >
        <text>进行中</text>
      </view>
      <view
        class="tab"
        :class="{ 'tab--active': activeTab === 'done' }"
        @tap="activeTab = 'done'"
      >
        <text>已完成</text>
      </view>
    </view>

    <!-- 课程列表 -->
    <view class="course-list">
      <view
        class="course-card"
        v-for="item in filteredCourses"
        :key="item.id"
      >
        <!-- 头部：课程名 + 类型标签 -->
        <view class="card-header">
          <view class="card-title-area">
            <text class="card-title">{{ item.name }}</text>
            <view
              class="type-tag"
              :class="{
                'type-tag--private': item.type === 'private',
                'type-tag--group': item.type === 'group',
              }"
            >
              <text class="type-tag-text">{{ typeLabel(item.type) }}</text>
            </view>
          </view>
        </view>

        <!-- 教练信息 -->
        <view class="coach-row">
          <view class="coach-avatar" :style="{ background: item.coachAvatarBg }">
            <text class="coach-avatar-text">{{ item.coachName.charAt(0) }}</text>
          </view>
          <text class="coach-name">教练：{{ item.coachName }}</text>
        </view>

        <!-- 进度条 -->
        <view class="progress-section">
          <view class="progress-info">
            <text class="progress-label">课程进度</text>
            <text class="progress-count">已上 {{ item.completed }}/{{ item.total }} 次</text>
          </view>
          <view class="progress-bar">
            <view
              class="progress-fill"
              :style="{ width: progressPercent(item) + '%' }"
              :class="{ 'progress-done': item.status === 'done' }"
            />
          </view>
        </view>

        <!-- 下次上课（仅进行中） -->
        <view class="next-class" v-if="item.status === 'active' && item.nextTime">
          <text class="next-label">下次上课</text>
          <text class="next-time">{{ item.nextTime }}</text>
          <text class="next-venue">{{ item.nextVenue }}</text>
        </view>

        <!-- 操作按钮 -->
        <view class="card-actions">
          <view class="action-btn action-btn--outline" @tap="onViewRecords(item.id)">
            <text class="action-btn-text action-btn-text--outline">查看记录</text>
          </view>
          <view
            v-if="item.status === 'done'"
            class="action-btn action-btn--primary"
            @tap="onRate(item.id)"
          >
            <text class="action-btn-text action-btn-text--primary">评价</text>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="empty" v-if="filteredCourses.length === 0">
        <text class="empty-text">暂无{{ activeTab === 'active' ? '进行中' : '已完成' }}的课程</text>
        <view class="empty-btn" @tap="goShop">
          <text class="empty-btn-text">去选课</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// ========== 类型 ==========
interface Course {
  readonly id: number
  readonly name: string
  readonly type: 'private' | 'group'
  readonly coachName: string
  readonly coachAvatarBg: string
  readonly completed: number
  readonly total: number
  readonly status: 'active' | 'done'
  readonly nextTime: string
  readonly nextVenue: string
}

// ========== Tab 切换 ==========
const activeTab = ref<'active' | 'done'>('active')

// ========== Mock 数据 ==========
const courses = ref<ReadonlyArray<Course>>([
  {
    id: 1,
    name: '正反手提升训练',
    type: 'private',
    coachName: '王教练',
    coachAvatarBg: '#2255CC',
    completed: 3,
    total: 10,
    status: 'active',
    nextTime: '03/25 周二 14:00-15:00',
    nextVenue: '1号室内场',
  },
  {
    id: 2,
    name: '发球专项训练',
    type: 'private',
    coachName: '张教练',
    coachAvatarBg: '#F59E0B',
    completed: 7,
    total: 10,
    status: 'active',
    nextTime: '03/26 周三 10:00-11:00',
    nextVenue: '2号室内场',
  },
  {
    id: 3,
    name: '初级网球团课',
    type: 'group',
    coachName: '李教练',
    coachAvatarBg: '#22C55E',
    completed: 4,
    total: 8,
    status: 'active',
    nextTime: '03/27 周四 16:00-17:30',
    nextVenue: '3号室外场',
  },
  {
    id: 4,
    name: '5次体验课包',
    type: 'private',
    coachName: '李教练',
    coachAvatarBg: '#22C55E',
    completed: 5,
    total: 5,
    status: 'done',
    nextTime: '',
    nextVenue: '',
  },
  {
    id: 5,
    name: '中级提高团课',
    type: 'group',
    coachName: '王教练',
    coachAvatarBg: '#2255CC',
    completed: 12,
    total: 12,
    status: 'done',
    nextTime: '',
    nextVenue: '',
  },
])

// ========== 筛选 ==========
const filteredCourses = computed(() => {
  return courses.value.filter((c) => c.status === activeTab.value)
})

// ========== 类型标签 ==========
function typeLabel(type: string): string {
  const map: Record<string, string> = {
    private: '私教',
    group: '团课',
  }
  return map[type] || type
}

// ========== 进度 ==========
function progressPercent(item: Course): number {
  if (item.total === 0) return 0
  return Math.min((item.completed / item.total) * 100, 100)
}

// ========== 操作 ==========
function onViewRecords(id: number) {
  uni.showToast({ title: '查看课程记录', icon: 'none' })
}

function onRate(id: number) {
  uni.showToast({ title: '打开评价', icon: 'none' })
}

function goShop() {
  uni.switchTab({ url: '/pages/shop/index' })
}
</script>

<style lang="scss" scoped>
.my-courses {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 120rpx;
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

// ========== 课程列表 ==========
.course-list {
  padding: $spacing-sm $spacing-base;
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

// ========== 课程卡片 ==========
.course-card {
  background: $color-bg-card;
  border-radius: $radius-lg;
  padding: $spacing-base;
}

.card-header {
  margin-bottom: $spacing-sm;
}

.card-title-area {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.card-title {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-text;
}

// ========== 类型标签 ==========
.type-tag {
  padding: 2rpx 14rpx;
  border-radius: $radius-sm;
  flex-shrink: 0;
}

.type-tag-text {
  font-size: $font-size-xs;
  font-weight: 600;
}

.type-tag--private {
  background: #EEF2FF;
  .type-tag-text { color: $brand-primary; }
}

.type-tag--group {
  background: #ECFDF5;
  .type-tag-text { color: $color-success; }
}

// ========== 教练信息 ==========
.coach-row {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-sm;
}

.coach-avatar {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.coach-avatar-text {
  color: #fff;
  font-size: $font-size-xs;
  font-weight: 700;
}

.coach-name {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

// ========== 进度条 ==========
.progress-section {
  margin-bottom: $spacing-sm;
}

.progress-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-xs;
}

.progress-label {
  font-size: $font-size-xs;
  color: $color-text-secondary;
}

.progress-count {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

.progress-bar {
  height: 16rpx;
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

.progress-done {
  background: $color-success;
}

// ========== 下次上课 ==========
.next-class {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  padding: $spacing-sm 0;
  border-top: 1rpx solid $color-border;
  margin-bottom: $spacing-sm;
}

.next-label {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
  flex-shrink: 0;
}

.next-time {
  font-size: $font-size-sm;
  color: $brand-primary;
  font-weight: 600;
}

.next-venue {
  font-size: $font-size-xs;
  color: $color-text-secondary;
  flex-shrink: 0;
}

// ========== 操作按钮 ==========
.card-actions {
  display: flex;
  gap: $spacing-sm;
  padding-top: $spacing-sm;
  border-top: 1rpx solid $color-border;
}

.action-btn {
  flex: 1;
  text-align: center;
  padding: $spacing-xs 0;
  border-radius: $radius-round;
}

.action-btn--outline {
  border: 1rpx solid $color-border;
}

.action-btn--primary {
  background: $brand-primary;
}

.action-btn-text {
  font-size: $font-size-sm;
  font-weight: 600;
}

.action-btn-text--outline {
  color: $color-text-secondary;
}

.action-btn-text--primary {
  color: #fff;
}

// ========== 空状态 ==========
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xl 0;
  gap: $spacing-base;
}

.empty-text {
  font-size: $font-size-base;
  color: $color-text-placeholder;
}

.empty-btn {
  background: $brand-primary;
  padding: $spacing-sm $spacing-lg;
  border-radius: $radius-round;
}

.empty-btn-text {
  color: #fff;
  font-size: $font-size-sm;
  font-weight: 600;
}
</style>
