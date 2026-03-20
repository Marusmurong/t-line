<template>
  <div class="students-page">
    <!-- 顶部标题 -->
    <div class="page-header">
      <h2 class="page-title">学员管理</h2>
    </div>

    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <a-row :gutter="16" align="center">
        <a-col :span="8">
          <a-input
            v-model="searchKeyword"
            placeholder="搜索学员姓名/手机号..."
            allow-clear
          >
            <template #prefix><icon-search /></template>
          </a-input>
        </a-col>
        <a-col :span="4">
          <a-select
            v-model="filterStatus"
            placeholder="学员状态"
            allow-clear
            style="width: 100%"
          >
            <a-option value="normal">正常</a-option>
            <a-option value="paused">暂停</a-option>
            <a-option value="lost">流失</a-option>
          </a-select>
        </a-col>
      </a-row>
    </a-card>

    <!-- 学员表格 -->
    <a-card style="margin-top: 16px">
      <a-table
        :data="filteredStudents"
        :pagination="pagination"
        stripe
        @page-change="onPageChange"
      >
        <template #columns>
          <a-table-column title="学员姓名" :width="140">
            <template #cell="{ record }">
              <div class="student-cell">
                <a-avatar :size="28" :style="{ background: record.avatarColor }">
                  {{ record.name.charAt(0) }}
                </a-avatar>
                <span class="student-name">{{ record.name }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="手机号" data-index="phone" :width="140" />
          <a-table-column title="责任教练" data-index="coach" :width="120" />
          <a-table-column title="累计课时" :width="100">
            <template #cell="{ record }">
              {{ record.totalHours }} 节
            </template>
          </a-table-column>
          <a-table-column title="最近上课" data-index="lastClass" :width="140" />
          <a-table-column title="状态" :width="100">
            <template #cell="{ record }">
              <a-tag :color="statusColorMap[record.status]">
                {{ statusLabelMap[record.status] }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="120">
            <template #cell="{ record }">
              <a-button type="text" size="small" @click="onViewRecords(record.id)">
                查看记录
              </a-button>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { IconSearch } from '@arco-design/web-vue/es/icon'

// ========== 类型 ==========
interface Student {
  readonly id: number
  readonly name: string
  readonly avatarColor: string
  readonly phone: string
  readonly coach: string
  readonly totalHours: number
  readonly lastClass: string
  readonly status: 'normal' | 'paused' | 'lost'
}

// ========== 颜色/标签映射 ==========
const statusColorMap: Record<string, string> = {
  normal: 'green',
  paused: 'gold',
  lost: 'red',
}

const statusLabelMap: Record<string, string> = {
  normal: '正常',
  paused: '暂停',
  lost: '流失',
}

// ========== 筛选状态 ==========
const searchKeyword = ref('')
const filterStatus = ref<string | undefined>(undefined)

// ========== 分页 ==========
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 8,
})

function onPageChange(page: number) {
  pagination.current = page
}

// ========== Mock 数据 ==========
const students = ref<ReadonlyArray<Student>>([
  { id: 1, name: '张伟', avatarColor: '#2255CC', phone: '138****1234', coach: '王教练', totalHours: 32, lastClass: '2026-03-19', status: 'normal' },
  { id: 2, name: '李娜', avatarColor: '#22C55E', phone: '139****5678', coach: '张教练', totalHours: 28, lastClass: '2026-03-18', status: 'normal' },
  { id: 3, name: '王明', avatarColor: '#F59E0B', phone: '137****9012', coach: '王教练', totalHours: 45, lastClass: '2026-03-20', status: 'normal' },
  { id: 4, name: '刘洋', avatarColor: '#8B5CF6', phone: '136****3456', coach: '李教练', totalHours: 18, lastClass: '2026-03-15', status: 'normal' },
  { id: 5, name: '陈静', avatarColor: '#EF4444', phone: '135****7890', coach: '赵教练', totalHours: 12, lastClass: '2026-02-28', status: 'paused' },
  { id: 6, name: '赵磊', avatarColor: '#06B6D4', phone: '133****2345', coach: '张教练', totalHours: 56, lastClass: '2026-03-17', status: 'normal' },
  { id: 7, name: '孙芳', avatarColor: '#EC4899', phone: '132****6789', coach: '李教练', totalHours: 8, lastClass: '2026-01-10', status: 'lost' },
  { id: 8, name: '周杰', avatarColor: '#14B8A6', phone: '131****0123', coach: '王教练', totalHours: 22, lastClass: '2026-03-12', status: 'paused' },
])

// ========== 筛选 ==========
const filteredStudents = computed(() => {
  return students.value.filter((s) => {
    if (filterStatus.value && s.status !== filterStatus.value) return false
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase()
      if (!s.name.toLowerCase().includes(kw) && !s.phone.includes(kw)) return false
    }
    return true
  })
})

// ========== 操作 ==========
function onViewRecords(id: number) {
  // TODO: 接入真实路由导航
  console.log('查看学员记录:', id)
}
</script>

<style lang="scss" scoped>
.students-page {
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

  .filter-card {
    :deep(.arco-card-body) {
      padding: 16px;
    }
  }

  .student-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .student-name {
    font-size: 14px;
    color: #1d2129;
  }
}
</style>
