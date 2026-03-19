<template>
  <div class="activity-list">
    <!-- 顶部标题栏 -->
    <div class="page-header">
      <h2 class="page-title">活动管理</h2>
      <a-button type="primary">
        <template #icon><icon-plus /></template>
        创建活动
      </a-button>
    </div>

    <!-- 活动卡片网格 -->
    <a-row :gutter="16">
      <a-col
        v-for="activity in activities"
        :key="activity.id"
        :span="8"
        style="margin-bottom: 16px"
      >
        <a-card class="activity-card" hoverable>
          <!-- 封面色块 + 状态角标 -->
          <template #cover>
            <div
              class="cover-block"
              :style="{ background: activity.coverColor }"
            >
              <a-tag
                class="status-badge"
                :color="statusColorMap[activity.status]"
                size="small"
              >
                {{ activity.status }}
              </a-tag>
            </div>
          </template>

          <!-- 活动信息 -->
          <div class="activity-info">
            <div class="activity-title-row">
              <span class="activity-name">{{ activity.name }}</span>
              <a-tag size="small" :color="typeColorMap[activity.type]">
                {{ activity.type }}
              </a-tag>
            </div>

            <div class="activity-meta">
              <div class="meta-item">
                <icon-clock-circle />
                <span>{{ activity.time }}</span>
              </div>
              <div class="meta-item">
                <icon-location />
                <span>{{ activity.venue }}</span>
              </div>
            </div>

            <!-- 报名进度 -->
            <div class="enrollment-section">
              <div class="enrollment-text">
                报名 {{ activity.enrolled }}/{{ activity.capacity }} 人
              </div>
              <a-progress
                :percent="activity.enrolled / activity.capacity"
                :color="activity.enrolled >= activity.capacity ? '#22C55E' : '#2255CC'"
                size="small"
              />
            </div>

            <!-- 操作按钮 -->
            <div class="activity-actions">
              <a-button type="text" size="small">编辑</a-button>
              <a-button
                v-if="activity.status !== '已结束'"
                type="text"
                size="small"
                status="danger"
              >
                取消
              </a-button>
              <a-button type="text" size="small">查看报名</a-button>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { IconPlus, IconClockCircle, IconLocation } from '@arco-design/web-vue/es/icon'

const statusColorMap: Record<string, string> = {
  '进行中': 'green',
  '草稿': 'gray',
  '已结束': 'red',
}

const typeColorMap: Record<string, string> = {
  '比赛': 'blue',
  '训练营': 'orange',
  '社交': 'purple',
  '公开课': 'green',
}

interface Activity {
  readonly id: number
  readonly name: string
  readonly type: string
  readonly status: string
  readonly time: string
  readonly venue: string
  readonly enrolled: number
  readonly capacity: number
  readonly coverColor: string
}

const activities = ref<ReadonlyArray<Activity>>([
  {
    id: 1,
    name: '春季网球公开赛',
    type: '比赛',
    status: '进行中',
    time: '2026-03-22 09:00 - 17:00',
    venue: '室内馆1号场',
    enrolled: 28,
    capacity: 32,
    coverColor: 'linear-gradient(135deg, #667eea, #764ba2)',
  },
  {
    id: 2,
    name: '周末双打友谊赛',
    type: '社交',
    status: '进行中',
    time: '2026-03-21 14:00 - 18:00',
    venue: '网球场A',
    enrolled: 12,
    capacity: 16,
    coverColor: 'linear-gradient(135deg, #f093fb, #f5576c)',
  },
  {
    id: 3,
    name: '青少年暑期训练营',
    type: '训练营',
    status: '草稿',
    time: '2026-07-01 - 2026-07-15',
    venue: '室内馆全场',
    enrolled: 0,
    capacity: 24,
    coverColor: 'linear-gradient(135deg, #4facfe, #00f2fe)',
  },
  {
    id: 4,
    name: '网球技术公开课',
    type: '公开课',
    status: '草稿',
    time: '2026-04-05 10:00 - 12:00',
    venue: '训练场C',
    enrolled: 0,
    capacity: 40,
    coverColor: 'linear-gradient(135deg, #43e97b, #38f9d7)',
  },
  {
    id: 5,
    name: '元旦迎新赛',
    type: '比赛',
    status: '已结束',
    time: '2026-01-01 09:00 - 17:00',
    venue: '网球场A+B',
    enrolled: 32,
    capacity: 32,
    coverColor: 'linear-gradient(135deg, #a18cd1, #fbc2eb)',
  },
  {
    id: 6,
    name: '冬季体能训练营',
    type: '训练营',
    status: '已结束',
    time: '2026-01-15 - 2026-02-15',
    venue: '室内馆2号场',
    enrolled: 18,
    capacity: 20,
    coverColor: 'linear-gradient(135deg, #ffecd2, #fcb69f)',
  },
])
</script>

<style lang="scss" scoped>
.activity-list {
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

  .activity-card {
    :deep(.arco-card-body) {
      padding: 0;
    }

    .cover-block {
      position: relative;
      height: 120px;
      border-radius: 4px 4px 0 0;
    }

    .status-badge {
      position: absolute;
      top: 10px;
      right: 10px;
    }
  }

  .activity-info {
    padding: 16px;
  }

  .activity-title-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
  }

  .activity-name {
    font-size: 16px;
    font-weight: 600;
    color: #1d2129;
  }

  .activity-meta {
    margin-bottom: 12px;

    .meta-item {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: #666;
      margin-bottom: 4px;

      :deep(.arco-icon) {
        font-size: 14px;
        color: #999;
      }
    }
  }

  .enrollment-section {
    margin-bottom: 12px;

    .enrollment-text {
      font-size: 13px;
      color: #666;
      margin-bottom: 4px;
    }
  }

  .activity-actions {
    display: flex;
    gap: 4px;
    border-top: 1px solid #f0f0f0;
    padding-top: 10px;
  }
}
</style>
