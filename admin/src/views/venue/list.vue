<template>
  <div class="venue-list">
    <!-- 顶部标题栏 -->
    <div class="page-header">
      <h2 class="page-title">场地管理</h2>
      <a-button type="primary">
        <template #icon><icon-plus /></template>
        新增场地
      </a-button>
    </div>

    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <a-row :gutter="16" align="center">
        <a-col :span="6">
          <a-select
            v-model="filterType"
            placeholder="场地类型"
            allow-clear
            style="width: 100%"
          >
            <a-option value="indoor">室内</a-option>
            <a-option value="outdoor">室外</a-option>
          </a-select>
        </a-col>
        <a-col :span="6">
          <a-select
            v-model="filterStatus"
            placeholder="状态"
            allow-clear
            style="width: 100%"
          >
            <a-option value="active">正常</a-option>
            <a-option value="maintenance">维护</a-option>
            <a-option value="disabled">停用</a-option>
          </a-select>
        </a-col>
        <a-col :span="8">
          <a-input
            v-model="searchKeyword"
            placeholder="搜索场地名称..."
            allow-clear
          >
            <template #prefix><icon-search /></template>
          </a-input>
        </a-col>
      </a-row>
    </a-card>

    <!-- 场地表格 -->
    <a-card style="margin-top: 16px">
      <a-table :data="filteredVenues" :pagination="false" stripe>
        <template #columns>
          <a-table-column title="场地名称" data-index="name">
            <template #cell="{ record }">
              <div class="venue-name">
                <span class="name-text">{{ record.name }}</span>
                <span class="name-sub">{{ record.description }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="类型" data-index="type" :width="100">
            <template #cell="{ record }">
              <a-tag :color="record.type === '室内' ? 'blue' : 'green'">
                {{ record.type }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="状态" data-index="status" :width="100">
            <template #cell="{ record }">
              <a-tag :color="statusColorMap[record.status]">
                {{ record.status }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="价格/小时" data-index="price" :width="120">
            <template #cell="{ record }">
              <span class="price-text">¥{{ record.price }}</span>
            </template>
          </a-table-column>
          <a-table-column title="今日预订率" data-index="bookingRate" :width="180">
            <template #cell="{ record }">
              <a-progress
                :percent="record.bookingRate / 100"
                :color="record.bookingRate > 80 ? '#22C55E' : record.bookingRate > 50 ? '#2255CC' : '#F59E0B'"
                size="small"
                style="width: 120px"
              />
              <span class="rate-text">{{ record.bookingRate }}%</span>
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="160">
            <template #cell="{ record }">
              <a-button type="text" size="small">编辑</a-button>
              <a-button
                type="text"
                size="small"
                :status="record.status === '停用' ? 'success' : 'danger'"
              >
                {{ record.status === '停用' ? '启用' : '停用' }}
              </a-button>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <a-pagination
          v-model:current="currentPage"
          :total="venues.length"
          :page-size="10"
          show-total
        />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { IconPlus, IconSearch } from '@arco-design/web-vue/es/icon'

const filterType = ref<string | undefined>(undefined)
const filterStatus = ref<string | undefined>(undefined)
const searchKeyword = ref('')
const currentPage = ref(1)

const statusColorMap: Record<string, string> = {
  '正常': 'green',
  '维护': 'orangered',
  '停用': 'red',
}

const typeFilterMap: Record<string, string> = {
  indoor: '室内',
  outdoor: '室外',
}

const statusFilterMap: Record<string, string> = {
  active: '正常',
  maintenance: '维护',
  disabled: '停用',
}

const venues = ref([
  { id: 1, name: '网球场A', description: '标准硬地球场', type: '室外', status: '正常', price: 200, bookingRate: 85 },
  { id: 2, name: '网球场B', description: '标准硬地球场', type: '室外', status: '正常', price: 200, bookingRate: 72 },
  { id: 3, name: '室内馆1号场', description: '恒温室内球场', type: '室内', status: '正常', price: 320, bookingRate: 92 },
  { id: 4, name: '室内馆2号场', description: '恒温室内球场', type: '室内', status: '维护', price: 320, bookingRate: 0 },
  { id: 5, name: '训练场C', description: '训练专用场地', type: '室外', status: '正常', price: 150, bookingRate: 45 },
  { id: 6, name: 'VIP球场', description: '高端独立球场', type: '室内', status: '停用', price: 500, bookingRate: 0 },
])

const filteredVenues = computed(() => {
  return venues.value.filter((v) => {
    if (filterType.value && v.type !== typeFilterMap[filterType.value]) return false
    if (filterStatus.value && v.status !== statusFilterMap[filterStatus.value]) return false
    if (searchKeyword.value && !v.name.includes(searchKeyword.value)) return false
    return true
  })
})
</script>

<style lang="scss" scoped>
.venue-list {
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

  .venue-name {
    display: flex;
    flex-direction: column;

    .name-text {
      font-weight: 500;
      color: #1d2129;
    }

    .name-sub {
      font-size: 12px;
      color: #999;
      margin-top: 2px;
    }
  }

  .price-text {
    font-weight: 600;
    color: #1d2129;
  }

  .rate-text {
    margin-left: 8px;
    font-size: 12px;
    color: #666;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }
}
</style>
