<template>
  <div class="coaches-page">
    <!-- 顶部标题栏 -->
    <div class="page-header">
      <h2 class="page-title">教练管理</h2>
      <a-button type="primary">
        <template #icon><icon-plus /></template>
        添加教练
      </a-button>
    </div>

    <!-- 教练卡片网格 -->
    <a-row :gutter="16">
      <a-col
        v-for="coach in coaches"
        :key="coach.id"
        :span="6"
        style="margin-bottom: 16px"
      >
        <a-card class="coach-card" hoverable>
          <div class="coach-info">
            <!-- 头像 -->
            <div class="avatar-wrapper">
              <div class="coach-avatar" :style="{ background: coach.avatarBg }">
                <span class="avatar-text">{{ coach.name.charAt(0) }}</span>
              </div>
              <span
                class="status-dot"
                :class="{
                  'status-dot--active': coach.status === 'active',
                  'status-dot--leave': coach.status === 'leave',
                }"
              ></span>
            </div>

            <!-- 姓名 + 等级 -->
            <div class="name-row">
              <span class="coach-name">{{ coach.name }}</span>
              <a-tag size="small" :color="levelColorMap[coach.level]">
                {{ coach.level }}
              </a-tag>
            </div>

            <!-- 课时 -->
            <div class="stat-row">
              <span class="stat-label">本月课时</span>
              <span class="stat-value">{{ coach.monthlyHours }} 节</span>
            </div>

            <!-- 评分 -->
            <div class="stat-row">
              <span class="stat-label">学员评分</span>
              <span class="rating-stars">
                <span
                  v-for="s in 5"
                  :key="s"
                  class="star"
                  :class="{ 'star--filled': s <= Math.round(coach.rating) }"
                >&#9733;</span>
              </span>
              <span class="rating-num">{{ coach.rating.toFixed(1) }}</span>
            </div>

            <!-- 状态 -->
            <div class="stat-row">
              <span class="stat-label">状态</span>
              <a-tag
                size="small"
                :color="coach.status === 'active' ? 'green' : 'gold'"
              >
                {{ statusLabel(coach.status) }}
              </a-tag>
            </div>

            <!-- 操作 -->
            <div class="coach-actions">
              <a-button type="text" size="small">查看详情</a-button>
              <a-button type="text" size="small">排班</a-button>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 绩效排行 -->
    <a-card class="ranking-card" style="margin-top: 8px">
      <template #title>绩效排行（本月 Top 5）</template>
      <a-table :data="topCoaches" :pagination="false" stripe>
        <template #columns>
          <a-table-column title="排名" :width="80">
            <template #cell="{ rowIndex }">
              <span
                class="rank-badge"
                :class="{
                  'rank-gold': rowIndex === 0,
                  'rank-silver': rowIndex === 1,
                  'rank-bronze': rowIndex === 2,
                }"
              >
                {{ rowIndex + 1 }}
              </span>
            </template>
          </a-table-column>
          <a-table-column title="教练" :width="140">
            <template #cell="{ record }">
              <div class="table-coach-cell">
                <a-avatar :size="28" :style="{ background: record.avatarBg }">
                  {{ record.name.charAt(0) }}
                </a-avatar>
                <span>{{ record.name }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="课时" data-index="monthlyHours" :width="100">
            <template #cell="{ record }">
              {{ record.monthlyHours }} 节
            </template>
          </a-table-column>
          <a-table-column title="学员数" data-index="studentCount" :width="100" />
          <a-table-column title="评分" :width="100">
            <template #cell="{ record }">
              <span class="rating-num">{{ record.rating.toFixed(1) }}</span>
              <span class="star star--filled">&#9733;</span>
            </template>
          </a-table-column>
          <a-table-column title="收入" :width="120">
            <template #cell="{ record }">
              <span class="amount-text">¥{{ record.revenue.toLocaleString() }}</span>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { IconPlus } from '@arco-design/web-vue/es/icon'

// ========== 类型 ==========
interface Coach {
  readonly id: number
  readonly name: string
  readonly avatarBg: string
  readonly level: '初级' | '中级' | '高级' | '资深'
  readonly monthlyHours: number
  readonly rating: number
  readonly status: 'active' | 'leave'
  readonly studentCount: number
  readonly revenue: number
}

// ========== 颜色映射 ==========
const levelColorMap: Record<string, string> = {
  '初级': 'blue',
  '中级': 'green',
  '高级': 'orange',
  '资深': 'purple',
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    active: '在职',
    leave: '休假',
  }
  return map[status] || status
}

// ========== Mock 数据 ==========
const coaches = ref<ReadonlyArray<Coach>>([
  { id: 1, name: '王教练', avatarBg: '#2255CC', level: '资深', monthlyHours: 48, rating: 4.9, status: 'active', studentCount: 15, revenue: 28800 },
  { id: 2, name: '张教练', avatarBg: '#F59E0B', level: '高级', monthlyHours: 42, rating: 4.8, status: 'active', studentCount: 12, revenue: 23100 },
  { id: 3, name: '李教练', avatarBg: '#22C55E', level: '中级', monthlyHours: 36, rating: 4.6, status: 'active', studentCount: 10, revenue: 16200 },
  { id: 4, name: '赵教练', avatarBg: '#8B5CF6', level: '高级', monthlyHours: 38, rating: 4.7, status: 'active', studentCount: 11, revenue: 20900 },
  { id: 5, name: '陈教练', avatarBg: '#EC4899', level: '初级', monthlyHours: 24, rating: 4.3, status: 'active', studentCount: 6, revenue: 8400 },
  { id: 6, name: '刘教练', avatarBg: '#14B8A6', level: '中级', monthlyHours: 30, rating: 4.5, status: 'leave', studentCount: 8, revenue: 13500 },
  { id: 7, name: '孙教练', avatarBg: '#EF4444', level: '资深', monthlyHours: 44, rating: 4.8, status: 'active', studentCount: 14, revenue: 26400 },
  { id: 8, name: '周教练', avatarBg: '#06B6D4', level: '初级', monthlyHours: 20, rating: 4.2, status: 'active', studentCount: 5, revenue: 7000 },
])

// ========== 绩效排行（按课时降序，前5名） ==========
const topCoaches = computed(() => {
  return [...coaches.value]
    .sort((a, b) => b.monthlyHours - a.monthlyHours)
    .slice(0, 5)
})
</script>

<style lang="scss" scoped>
.coaches-page {
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

  .coach-card {
    :deep(.arco-card-body) {
      padding: 16px;
    }
  }

  .coach-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
  }

  // ========== 头像 ==========
  .avatar-wrapper {
    position: relative;
  }

  .coach-avatar {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .avatar-text {
    color: #fff;
    font-size: 24px;
    font-weight: 700;
  }

  .status-dot {
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    border: 2px solid #fff;
  }

  .status-dot--active {
    background: #22c55e;
  }

  .status-dot--leave {
    background: #f59e0b;
  }

  // ========== 名字行 ==========
  .name-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .coach-name {
    font-size: 16px;
    font-weight: 600;
    color: #1d2129;
  }

  // ========== 统计行 ==========
  .stat-row {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    justify-content: center;
  }

  .stat-label {
    font-size: 13px;
    color: #999;
  }

  .stat-value {
    font-size: 14px;
    font-weight: 600;
    color: #1d2129;
  }

  // ========== 评分星 ==========
  .rating-stars {
    display: flex;
    gap: 1px;
  }

  .star {
    font-size: 14px;
    color: #e0e0e0;
  }

  .star--filled {
    color: #f59e0b;
  }

  .rating-num {
    font-size: 13px;
    font-weight: 600;
    color: #f59e0b;
  }

  // ========== 操作 ==========
  .coach-actions {
    display: flex;
    gap: 4px;
    border-top: 1px solid #f0f0f0;
    padding-top: 10px;
    width: 100%;
    justify-content: center;
  }

  // ========== 排行表格 ==========
  .ranking-card {
    :deep(.arco-card-header-title) {
      font-size: 16px;
      font-weight: 600;
    }
  }

  .table-coach-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .rank-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 50%;
    font-size: 13px;
    font-weight: 700;
    background: #f0f0f0;
    color: #666;
  }

  .rank-gold {
    background: linear-gradient(135deg, #fbbf24, #f59e0b);
    color: #fff;
  }

  .rank-silver {
    background: linear-gradient(135deg, #d1d5db, #9ca3af);
    color: #fff;
  }

  .rank-bronze {
    background: linear-gradient(135deg, #f59e0b, #d97706);
    color: #fff;
  }

  .amount-text {
    font-weight: 600;
    color: #1d2129;
  }
}
</style>
