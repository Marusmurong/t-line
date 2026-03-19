<template>
  <view class="my-courses">
    <!-- 课程卡片列表 -->
    <view
      class="course-card"
      v-for="item in courses"
      :key="item.id"
    >
      <!-- 头部：课程名 + 状态 -->
      <view class="card-header">
        <view class="card-title-area">
          <text class="card-title">{{ item.name }}</text>
          <text class="card-coach">{{ item.coach }}</text>
        </view>
        <view
          class="status-tag"
          :class="{
            'status-tag--active': item.status === 'active',
            'status-tag--done': item.status === 'done',
          }"
        >
          <text class="status-tag-text">{{ statusLabel(item.status) }}</text>
        </view>
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

      <!-- 下次上课时间 -->
      <view class="next-class" v-if="item.status === 'active' && item.nextTime">
        <text class="next-label">下次上课</text>
        <text class="next-time">{{ item.nextTime }}</text>
      </view>
    </view>

    <!-- 空状态 -->
    <view class="empty" v-if="courses.length === 0">
      <text class="empty-text">暂无课程</text>
      <view class="empty-btn" @tap="goShop">
        <text class="empty-btn-text">去选课</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

// ========== 类型 ==========
interface Course {
  id: number
  name: string
  coach: string
  completed: number
  total: number
  status: 'active' | 'done'
  nextTime: string
}

// ========== Mock 数据 ==========
const courses = ref<Course[]>([
  {
    id: 1,
    name: '私教课 - 正反手提升',
    coach: '王教练',
    completed: 3,
    total: 10,
    status: 'active',
    nextTime: '03/25 周二 14:00-15:00',
  },
  {
    id: 2,
    name: '私教课 - 发球专项训练',
    coach: '张教练',
    completed: 7,
    total: 10,
    status: 'active',
    nextTime: '03/26 周三 10:00-11:00',
  },
  {
    id: 3,
    name: '5次体验课包',
    coach: '李教练',
    completed: 5,
    total: 5,
    status: 'done',
    nextTime: '',
  },
  {
    id: 4,
    name: '20次私教课包',
    coach: '王教练',
    completed: 20,
    total: 20,
    status: 'done',
    nextTime: '',
  },
])

// ========== 状态标签 ==========
function statusLabel(status: string): string {
  const map: Record<string, string> = {
    active: '进行中',
    done: '已完成',
  }
  return map[status] || status
}

// ========== 进度 ==========
function progressPercent(item: Course): number {
  if (item.total === 0) return 0
  return Math.min((item.completed / item.total) * 100, 100)
}

// ========== 导航 ==========
function goShop() {
  uni.switchTab({ url: '/pages/shop/index' })
}
</script>

<style lang="scss" scoped>
.my-courses {
  background: $color-bg;
  min-height: 100vh;
  padding: $spacing-sm $spacing-base;
  padding-bottom: 120rpx;
}

// ========== 课程卡片 ==========
.course-card {
  background: $color-bg-card;
  border-radius: $radius-lg;
  padding: $spacing-base;
  margin-bottom: $spacing-sm;
}

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: $spacing-sm;
}

.card-title-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.card-title {
  font-size: $font-size-base;
  font-weight: 700;
  color: $color-text;
}

.card-coach {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

// ========== 状态标签 ==========
.status-tag {
  padding: 4rpx 16rpx;
  border-radius: $radius-sm;
  flex-shrink: 0;
}

.status-tag-text {
  font-size: $font-size-xs;
  font-weight: 600;
}

.status-tag--active {
  background: #EEF2FF;
  .status-tag-text { color: $brand-primary; }
}

.status-tag--done {
  background: #ECFDF5;
  .status-tag-text { color: $color-success; }
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
  padding-top: $spacing-sm;
  border-top: 1rpx solid $color-border;
}

.next-label {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
}

.next-time {
  font-size: $font-size-sm;
  color: $brand-primary;
  font-weight: 600;
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
