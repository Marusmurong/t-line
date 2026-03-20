<template>
  <div class="schedules-page">
    <!-- 顶部标题 -->
    <div class="page-header">
      <h2 class="page-title">课程管理</h2>
      <a-button type="primary">
        <template #icon><icon-plus /></template>
        新建课程
      </a-button>
    </div>

    <div class="main-layout">
      <!-- 左侧：课程列表 -->
      <div class="left-panel">
        <a-card class="list-card">
          <!-- Tab -->
          <a-tabs v-model:active-key="activeType" type="capsule" size="small">
            <a-tab-pane key="all" title="全部" />
            <a-tab-pane key="private" title="私教" />
            <a-tab-pane key="group" title="团课" />
            <a-tab-pane key="camp" title="训练营" />
          </a-tabs>

          <!-- 课程条目 -->
          <div class="schedule-list">
            <div
              class="schedule-item"
              v-for="item in filteredSchedules"
              :key="item.id"
              :class="{ 'schedule-item--active': selectedId === item.id }"
              @click="selectedId = item.id"
            >
              <div class="item-header">
                <span class="item-coach">{{ item.coach }}</span>
                <a-tag
                  size="small"
                  :color="typeColorMap[item.type]"
                >
                  {{ typeLabelMap[item.type] }}
                </a-tag>
              </div>
              <div class="item-time">{{ item.day }} {{ item.timeRange }}</div>
              <div class="item-meta">
                <span>{{ item.venue }}</span>
                <span>{{ item.students }}人</span>
                <a-tag
                  size="small"
                  :color="statusColorMap[item.status]"
                >
                  {{ item.status }}
                </a-tag>
              </div>
            </div>
          </div>
        </a-card>
      </div>

      <!-- 右侧：周视图日历 -->
      <div class="right-panel">
        <a-card class="calendar-card">
          <!-- 日历头部 -->
          <div class="calendar-header">
            <a-button size="small" @click="prevWeek">
              <template #icon><icon-left /></template>
            </a-button>
            <span class="week-label">{{ weekLabel }}</span>
            <a-button size="small" @click="nextWeek">
              <template #icon><icon-right /></template>
            </a-button>
          </div>

          <!-- 日历网格 -->
          <div class="calendar-grid">
            <!-- 表头：周一到周日 -->
            <div class="grid-header">
              <div class="time-gutter"></div>
              <div
                class="day-header"
                v-for="day in weekDays"
                :key="day.label"
              >
                <span class="day-name">{{ day.label }}</span>
                <span class="day-date">{{ day.date }}</span>
              </div>
            </div>

            <!-- 时间行 -->
            <div class="grid-body">
              <div class="time-row" v-for="hour in hours" :key="hour">
                <div class="time-gutter">
                  <span class="time-text">{{ hour }}:00</span>
                </div>
                <div
                  class="day-cell"
                  v-for="(day, dayIdx) in weekDays"
                  :key="day.label"
                >
                  <!-- 课程块 -->
                  <a-tooltip
                    v-for="block in getBlocks(dayIdx, hour)"
                    :key="block.id"
                    :content="`${block.coach} · ${typeLabelMap[block.type]}\n${block.timeRange} · ${block.venue}\n学员: ${block.students}人`"
                  >
                    <div
                      class="course-block"
                      :class="`block-${block.type}`"
                      :style="{ height: block.height + 'px' }"
                    >
                      <span class="block-coach">{{ block.coach }}</span>
                      <span class="block-time">{{ block.timeRange }}</span>
                    </div>
                  </a-tooltip>
                </div>
              </div>
            </div>
          </div>

          <!-- 图例 -->
          <div class="legend">
            <span class="legend-item">
              <span class="legend-dot" style="background: #2255cc"></span>私教
            </span>
            <span class="legend-item">
              <span class="legend-dot" style="background: #22c55e"></span>团课
            </span>
            <span class="legend-item">
              <span class="legend-dot" style="background: #f59e0b"></span>训练营
            </span>
          </div>
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { IconPlus, IconLeft, IconRight } from '@arco-design/web-vue/es/icon'

// ========== 类型 ==========
interface Schedule {
  readonly id: number
  readonly coach: string
  readonly type: 'private' | 'group' | 'camp'
  readonly day: string
  readonly dayIndex: number // 0=周一 ... 6=周日
  readonly startHour: number
  readonly duration: number // 小时
  readonly timeRange: string
  readonly venue: string
  readonly students: number
  readonly status: string
}

// ========== 颜色映射 ==========
const typeColorMap: Record<string, string> = {
  private: 'blue',
  group: 'green',
  camp: 'orange',
}

const typeLabelMap: Record<string, string> = {
  private: '私教',
  group: '团课',
  camp: '训练营',
}

const statusColorMap: Record<string, string> = {
  '已排': 'blue',
  '进行中': 'green',
  '已结束': 'gray',
}

// ========== Tab ==========
const activeType = ref('all')
const selectedId = ref<number | null>(null)

// ========== 周导航 ==========
const weekOffset = ref(0)

function prevWeek() {
  weekOffset.value = weekOffset.value - 1
}

function nextWeek() {
  weekOffset.value = weekOffset.value + 1
}

function getMonday(offset: number): Date {
  const now = new Date(2026, 2, 20) // 2026-03-20 是周五
  const day = now.getDay()
  const diff = now.getDate() - day + (day === 0 ? -6 : 1) + offset * 7
  return new Date(now.getFullYear(), now.getMonth(), diff)
}

const weekDays = computed(() => {
  const monday = getMonday(weekOffset.value)
  const labels = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  return labels.map((label, i) => {
    const d = new Date(monday)
    d.setDate(monday.getDate() + i)
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    return { label, date: `${mm}/${dd}` }
  })
})

const weekLabel = computed(() => {
  const monday = getMonday(weekOffset.value)
  const sunday = new Date(monday)
  sunday.setDate(monday.getDate() + 6)
  const fmt = (d: Date) => `${d.getMonth() + 1}/${d.getDate()}`
  return `${fmt(monday)} - ${fmt(sunday)}`
})

// ========== 时间轴 ==========
const hours = Array.from({ length: 13 }, (_, i) => i + 8) // 8:00 - 20:00

// ========== Mock 数据 ==========
const schedules = ref<ReadonlyArray<Schedule>>([
  { id: 1, coach: '王教练', type: 'private', day: '周一', dayIndex: 0, startHour: 9, duration: 1, timeRange: '09:00-10:00', venue: '1号室内场', students: 1, status: '已排' },
  { id: 2, coach: '张教练', type: 'group', day: '周一', dayIndex: 0, startHour: 14, duration: 1.5, timeRange: '14:00-15:30', venue: '3号室外场', students: 8, status: '已排' },
  { id: 3, coach: '李教练', type: 'private', day: '周二', dayIndex: 1, startHour: 10, duration: 1, timeRange: '10:00-11:00', venue: '2号室内场', students: 1, status: '已排' },
  { id: 4, coach: '赵教练', type: 'camp', day: '周三', dayIndex: 2, startHour: 9, duration: 3, timeRange: '09:00-12:00', venue: '全场地', students: 12, status: '进行中' },
  { id: 5, coach: '王教练', type: 'private', day: '周四', dayIndex: 3, startHour: 16, duration: 1, timeRange: '16:00-17:00', venue: '1号室内场', students: 1, status: '已排' },
  { id: 6, coach: '孙教练', type: 'group', day: '周五', dayIndex: 4, startHour: 18, duration: 1.5, timeRange: '18:00-19:30', venue: '2号室内场', students: 6, status: '已排' },
])

// ========== 筛选 ==========
const filteredSchedules = computed(() => {
  if (activeType.value === 'all') return schedules.value
  return schedules.value.filter((s) => s.type === activeType.value)
})

// ========== 日历块 ==========
function getBlocks(dayIndex: number, hour: number): ReadonlyArray<Schedule & { readonly height: number }> {
  return schedules.value
    .filter((s) => s.dayIndex === dayIndex && s.startHour === hour)
    .map((s) => ({
      ...s,
      height: s.duration * 48,
    }))
}
</script>

<style lang="scss" scoped>
.schedules-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: #1d2129;
    margin: 0;
  }

  // ========== 主布局 ==========
  .main-layout {
    display: flex;
    gap: 16px;
  }

  .left-panel {
    width: 40%;
    flex-shrink: 0;
  }

  .right-panel {
    width: 60%;
    flex-shrink: 0;
  }

  // ========== 课程列表 ==========
  .list-card {
    :deep(.arco-card-body) {
      padding: 16px;
    }
  }

  .schedule-list {
    margin-top: 12px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    max-height: 600px;
    overflow-y: auto;
  }

  .schedule-item {
    padding: 12px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      border-color: #2255cc;
      background: #f8faff;
    }
  }

  .schedule-item--active {
    border-color: #2255cc;
    background: #eef2ff;
  }

  .item-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
  }

  .item-coach {
    font-size: 15px;
    font-weight: 600;
    color: #1d2129;
  }

  .item-time {
    font-size: 13px;
    color: #666;
    margin-bottom: 4px;
  }

  .item-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: #999;
  }

  // ========== 日历 ==========
  .calendar-card {
    :deep(.arco-card-body) {
      padding: 16px;
    }
  }

  .calendar-header {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16px;
    margin-bottom: 16px;
  }

  .week-label {
    font-size: 15px;
    font-weight: 600;
    color: #1d2129;
  }

  .calendar-grid {
    overflow-x: auto;
  }

  .grid-header {
    display: grid;
    grid-template-columns: 50px repeat(7, 1fr);
    border-bottom: 1px solid #f0f0f0;
    padding-bottom: 8px;
    margin-bottom: 4px;
  }

  .time-gutter {
    width: 50px;
    flex-shrink: 0;
  }

  .day-header {
    text-align: center;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .day-name {
    font-size: 13px;
    font-weight: 600;
    color: #1d2129;
  }

  .day-date {
    font-size: 11px;
    color: #999;
  }

  .grid-body {
    position: relative;
  }

  .time-row {
    display: grid;
    grid-template-columns: 50px repeat(7, 1fr);
    min-height: 48px;
    border-bottom: 1px solid #fafafa;
  }

  .time-text {
    font-size: 11px;
    color: #999;
    line-height: 48px;
  }

  .day-cell {
    position: relative;
    border-left: 1px solid #fafafa;
    min-height: 48px;
  }

  // ========== 课程块 ==========
  .course-block {
    position: absolute;
    top: 0;
    left: 2px;
    right: 2px;
    border-radius: 4px;
    padding: 4px 6px;
    overflow: hidden;
    z-index: 1;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    gap: 2px;

    &:hover {
      opacity: 0.85;
    }
  }

  .block-private {
    background: rgba(34, 85, 204, 0.15);
    border-left: 3px solid #2255cc;
  }

  .block-group {
    background: rgba(34, 197, 94, 0.15);
    border-left: 3px solid #22c55e;
  }

  .block-camp {
    background: rgba(245, 158, 11, 0.15);
    border-left: 3px solid #f59e0b;
  }

  .block-coach {
    font-size: 12px;
    font-weight: 600;
    color: #1d2129;
  }

  .block-time {
    font-size: 11px;
    color: #666;
  }

  // ========== 图例 ==========
  .legend {
    display: flex;
    gap: 20px;
    margin-top: 16px;
    font-size: 13px;
    color: #666;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .legend-dot {
    width: 12px;
    height: 12px;
    border-radius: 3px;
  }
}
</style>
