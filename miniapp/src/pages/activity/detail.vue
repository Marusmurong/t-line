<template>
  <view class="activity-detail">
    <!-- 封面 -->
    <view class="cover" :style="{ background: activity.coverBg }">
      <view class="cover-tag">
        <text class="cover-tag-text">{{ activity.type }}</text>
      </view>
    </view>

    <!-- 活动信息 -->
    <view class="info-card">
      <text class="activity-title">{{ activity.title }}</text>

      <view class="info-row">
        <text class="info-label">时间</text>
        <text class="info-value">{{ activity.time }}</text>
      </view>
      <view class="info-row">
        <text class="info-label">场地</text>
        <text class="info-value">{{ activity.venue }}</text>
      </view>
      <view class="info-row" v-if="activity.coach">
        <text class="info-label">教练</text>
        <text class="info-value">{{ activity.coach }}</text>
      </view>
    </view>

    <!-- 报名状态 -->
    <view class="status-card">
      <view class="status-header">
        <text class="status-title">报名状态</text>
        <text class="status-count">{{ activity.enrolled }}/{{ activity.capacity }}人</text>
      </view>
      <view class="progress-bar">
        <view
          class="progress-fill"
          :style="{ width: progressPercent + '%' }"
          :class="{ 'progress-full': isFull }"
        />
      </view>
      <text class="status-hint" v-if="isFull">名额已满</text>
      <text class="status-hint" v-else>还剩 {{ activity.capacity - activity.enrolled }} 个名额</text>
    </view>

    <!-- 活动说明 -->
    <view class="desc-card">
      <text class="desc-title">活动说明</text>
      <text class="desc-content">{{ activity.description }}</text>
    </view>

    <!-- 水平要求 + 价格 -->
    <view class="extra-card">
      <view class="extra-row">
        <text class="extra-label">水平要求</text>
        <view class="level-tag" :class="levelClass">
          <text class="level-text">{{ activity.level }}</text>
        </view>
      </view>
      <view class="extra-row">
        <text class="extra-label">费用</text>
        <text class="extra-price">¥{{ activity.price }}/人</text>
      </view>
    </view>

    <!-- 底部固定按钮 -->
    <view class="bottom-bar">
      <view
        class="action-btn"
        :class="{
          'action-btn--enrolled': isEnrolled,
          'action-btn--full': isFull && !isEnrolled,
        }"
        @tap="onAction"
      >
        <text class="action-btn-text" v-if="isEnrolled">已报名</text>
        <text class="action-btn-text" v-else-if="isFull">已满员</text>
        <text class="action-btn-text" v-else>立即报名 ¥{{ activity.price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// ========== 类型 ==========
interface ActivityDetail {
  id: number
  title: string
  type: string
  time: string
  venue: string
  coach: string
  level: string
  enrolled: number
  capacity: number
  price: number
  coverBg: string
  description: string
}

// ========== 页面参数 ==========
const activityId = ref(0)

// ========== Mock 数据 ==========
const activity = ref<ActivityDetail>({
  id: 1,
  title: '周六下午畅打',
  type: '畅打',
  time: '2026年3月22日 14:00-17:00',
  venue: '1号室内场',
  coach: '王教练',
  level: '不限',
  enrolled: 12,
  capacity: 16,
  price: 68,
  coverBg: 'linear-gradient(135deg, #2255CC, #4477DD)',
  description: '每周六下午的固定畅打活动，适合所有水平的球友参加。\n\n活动内容：\n- 自由对打，可单打或双打\n- 教练现场指导\n- 提供饮用水\n\n注意事项：\n- 请自备球拍（也可现场租赁）\n- 请提前15分钟到场签到\n- 取消报名请提前24小时',
})

const isEnrolled = ref(false)

// ========== 接收页面参数 ==========
function onLoad(options: Record<string, string>) {
  if (options.id) {
    activityId.value = Number(options.id)
    // TODO: 接入真实 API
    // const data = await activityApi.getActivityDetail(activityId.value)
  }
}

// #ifdef MP-WEIXIN
defineExpose({ onLoad })
// #endif

// ========== 计算属性 ==========
const progressPercent = computed(() =>
  Math.min((activity.value.enrolled / activity.value.capacity) * 100, 100),
)

const isFull = computed(() =>
  activity.value.enrolled >= activity.value.capacity,
)

const levelClass = computed(() => {
  const map: Record<string, string> = {
    '不限': 'level--any',
    '初级': 'level--beginner',
    '中级': 'level--intermediate',
    '高级': 'level--advanced',
  }
  return map[activity.value.level] || 'level--any'
})

// ========== 操作 ==========
function onAction() {
  if (isEnrolled.value) {
    uni.showModal({
      title: '取消报名',
      content: '确定要取消报名吗？',
      success: (res) => {
        if (res.confirm) {
          isEnrolled.value = false
          activity.value = {
            ...activity.value,
            enrolled: activity.value.enrolled - 1,
          }
          uni.showToast({ title: '已取消报名', icon: 'none' })
        }
      },
    })
    return
  }

  if (isFull.value) return

  isEnrolled.value = true
  activity.value = {
    ...activity.value,
    enrolled: activity.value.enrolled + 1,
  }
  uni.showToast({ title: '报名成功', icon: 'success' })
}
</script>

<style lang="scss" scoped>
.activity-detail {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 160rpx;
}

// ========== 封面 ==========
.cover {
  width: 100%;
  height: 400rpx;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  padding: $spacing-base;
  position: relative;
}

.cover-tag {
  background: rgba(255, 255, 255, 0.9);
  padding: $spacing-xs $spacing-base;
  border-radius: $radius-sm;
}

.cover-tag-text {
  font-size: $font-size-sm;
  font-weight: 700;
  color: $brand-primary;
}

// ========== 信息卡片 ==========
.info-card {
  background: $color-bg-card;
  margin: -$spacing-lg $spacing-base 0;
  border-radius: $radius-lg;
  padding: $spacing-base;
  position: relative;
  z-index: 1;
}

.activity-title {
  font-size: $font-size-xl;
  font-weight: 700;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-sm;
}

.info-row {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  padding: $spacing-xs 0;
}

.info-label {
  font-size: $font-size-sm;
  color: $color-text-placeholder;
  width: 80rpx;
  flex-shrink: 0;
}

.info-value {
  font-size: $font-size-sm;
  color: $color-text;
}

// ========== 报名状态 ==========
.status-card {
  background: $color-bg-card;
  margin: $spacing-sm $spacing-base;
  border-radius: $radius-lg;
  padding: $spacing-base;
}

.status-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-sm;
}

.status-title {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.status-count {
  font-size: $font-size-sm;
  color: $brand-primary;
  font-weight: 600;
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

.progress-full {
  background: $color-text-disabled;
}

.status-hint {
  font-size: $font-size-xs;
  color: $color-text-placeholder;
  margin-top: $spacing-xs;
}

// ========== 活动说明 ==========
.desc-card {
  background: $color-bg-card;
  margin: 0 $spacing-base;
  border-radius: $radius-lg;
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

// ========== 额外信息 ==========
.extra-card {
  background: $color-bg-card;
  margin: $spacing-sm $spacing-base;
  border-radius: $radius-lg;
  padding: $spacing-base;
}

.extra-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-xs 0;
}

.extra-label {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.extra-price {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-error;
}

// ========== 水平标签 ==========
.level-tag {
  padding: 4rpx 16rpx;
  border-radius: $radius-sm;
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

.action-btn {
  background: $brand-primary;
  padding: $spacing-sm 0;
  border-radius: $radius-round;
  text-align: center;
}

.action-btn-text {
  color: #fff;
  font-size: $font-size-base;
  font-weight: 700;
}

.action-btn--enrolled {
  background: $color-success;
}

.action-btn--full {
  background: $color-text-disabled;
}
</style>
