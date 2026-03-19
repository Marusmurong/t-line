<template>
  <div class="order-list">
    <!-- 状态 Tab -->
    <div class="page-header">
      <h2 class="page-title">订单管理</h2>
    </div>

    <a-card class="filter-card">
      <a-tabs v-model:active-key="activeStatus" type="capsule" size="small">
        <a-tab-pane key="all" title="全部" />
        <a-tab-pane key="pending" title="待支付" />
        <a-tab-pane key="paid" title="已支付" />
        <a-tab-pane key="completed" title="已完成" />
        <a-tab-pane key="cancelled" title="已取消" />
        <a-tab-pane key="refunding" title="退款中" />
      </a-tabs>

      <!-- 筛选行 -->
      <a-row :gutter="16" style="margin-top: 12px" align="center">
        <a-col :span="5">
          <a-select
            v-model="filterOrderType"
            placeholder="订单类型"
            allow-clear
            style="width: 100%"
          >
            <a-option value="venue">场地预约</a-option>
            <a-option value="course">课程购买</a-option>
            <a-option value="activity">活动报名</a-option>
            <a-option value="product">商品购买</a-option>
          </a-select>
        </a-col>
        <a-col :span="8">
          <a-range-picker
            v-model="dateRange"
            style="width: 100%"
            placeholder="['开始日期', '结束日期']"
          />
        </a-col>
        <a-col :span="6">
          <a-input
            v-model="searchKeyword"
            placeholder="搜索订单号/用户名..."
            allow-clear
          >
            <template #prefix><icon-search /></template>
          </a-input>
        </a-col>
      </a-row>
    </a-card>

    <!-- 订单表格 -->
    <a-card style="margin-top: 16px">
      <a-table :data="filteredOrders" :pagination="false" stripe>
        <template #columns>
          <a-table-column title="订单号" data-index="orderNo" :width="160" />
          <a-table-column title="用户" data-index="user" :width="140">
            <template #cell="{ record }">
              <div class="user-cell">
                <a-avatar :size="28" :style="{ background: record.avatarColor }">
                  {{ record.user.charAt(0) }}
                </a-avatar>
                <span class="user-name">{{ record.user }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="订单类型" data-index="type" :width="110">
            <template #cell="{ record }">
              <a-tag :color="orderTypeColor[record.type]">{{ record.type }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="金额" data-index="amount" :width="100">
            <template #cell="{ record }">
              <span class="amount-text">¥{{ record.amount }}</span>
            </template>
          </a-table-column>
          <a-table-column title="支付方式" data-index="payMethod" :width="100" />
          <a-table-column title="状态" data-index="status" :width="100">
            <template #cell="{ record }">
              <a-tag :color="orderStatusColor[record.status]">{{ record.status }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="创建时间" data-index="createdAt" :width="170" />
          <a-table-column title="操作" :width="140">
            <template #cell="{ record }">
              <a-button type="text" size="small">详情</a-button>
              <a-button
                v-if="record.status === '已支付' || record.status === '退款中'"
                type="text"
                size="small"
                status="warning"
              >
                退款
              </a-button>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <!-- 底部：分页 + 统计 -->
      <div class="table-footer">
        <div class="page-summary">
          本页总额：<span class="summary-amount">¥{{ pageTotalAmount }}</span>
        </div>
        <a-pagination
          v-model:current="currentPage"
          :total="orders.length"
          :page-size="10"
          show-total
        />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { IconSearch } from '@arco-design/web-vue/es/icon'

const activeStatus = ref('all')
const filterOrderType = ref<string | undefined>(undefined)
const dateRange = ref<string[] | undefined>(undefined)
const searchKeyword = ref('')
const currentPage = ref(1)

const orderTypeColor: Record<string, string> = {
  '场地预约': 'blue',
  '课程购买': 'green',
  '活动报名': 'orange',
  '商品购买': 'purple',
}

const orderStatusColor: Record<string, string> = {
  '待支付': 'orangered',
  '已支付': 'green',
  '已完成': 'blue',
  '已取消': 'red',
  '退款中': 'gold',
}

const statusFilterMap: Record<string, string> = {
  pending: '待支付',
  paid: '已支付',
  completed: '已完成',
  cancelled: '已取消',
  refunding: '退款中',
}

const typeFilterMap: Record<string, string> = {
  venue: '场地预约',
  course: '课程购买',
  activity: '活动报名',
  product: '商品购买',
}

const orders = ref([
  { orderNo: '2026032000001', user: '张伟', avatarColor: '#2255CC', type: '场地预约', amount: '280', payMethod: '微信支付', status: '已支付', createdAt: '2026-03-20 10:30:00' },
  { orderNo: '2026032000002', user: '李娜', avatarColor: '#22C55E', type: '课程购买', amount: '1280', payMethod: '支付宝', status: '已支付', createdAt: '2026-03-20 10:15:00' },
  { orderNo: '2026032000003', user: '王明', avatarColor: '#F59E0B', type: '活动报名', amount: '68', payMethod: '微信支付', status: '已完成', createdAt: '2026-03-20 09:45:00' },
  { orderNo: '2026032000004', user: '刘洋', avatarColor: '#8B5CF6', type: '商品购买', amount: '350', payMethod: '-', status: '待支付', createdAt: '2026-03-20 09:30:00' },
  { orderNo: '2026032000005', user: '陈静', avatarColor: '#EF4444', type: '场地预约', amount: '140', payMethod: '微信支付', status: '已取消', createdAt: '2026-03-20 09:10:00' },
  { orderNo: '2026032000006', user: '赵磊', avatarColor: '#06B6D4', type: '课程购买', amount: '2560', payMethod: '支付宝', status: '退款中', createdAt: '2026-03-20 08:50:00' },
  { orderNo: '2026032000007', user: '孙芳', avatarColor: '#EC4899', type: '场地预约', amount: '320', payMethod: '微信支付', status: '已支付', createdAt: '2026-03-19 18:20:00' },
  { orderNo: '2026032000008', user: '周杰', avatarColor: '#14B8A6', type: '活动报名', amount: '128', payMethod: '支付宝', status: '已完成', createdAt: '2026-03-19 17:45:00' },
])

const filteredOrders = computed(() => {
  return orders.value.filter((o) => {
    if (activeStatus.value !== 'all' && o.status !== statusFilterMap[activeStatus.value]) return false
    if (filterOrderType.value && o.type !== typeFilterMap[filterOrderType.value]) return false
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase()
      if (!o.orderNo.includes(kw) && !o.user.toLowerCase().includes(kw)) return false
    }
    return true
  })
})

const pageTotalAmount = computed(() => {
  const total = filteredOrders.value.reduce((sum, o) => sum + Number(o.amount), 0)
  return total.toLocaleString()
})
</script>

<style lang="scss" scoped>
.order-list {
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

  .user-cell {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .user-name {
    font-size: 14px;
    color: #1d2129;
  }

  .amount-text {
    font-weight: 600;
    color: #1d2129;
  }

  .table-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
  }

  .page-summary {
    font-size: 14px;
    color: #666;
  }

  .summary-amount {
    font-weight: 700;
    font-size: 16px;
    color: #2255CC;
  }
}
</style>
