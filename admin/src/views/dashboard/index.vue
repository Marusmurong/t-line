<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <a-row :gutter="16">
      <a-col :span="6" v-for="stat in stats" :key="stat.title">
        <a-card class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-title">{{ stat.title }}</div>
              <div class="stat-value" :style="{ color: stat.color }">{{ stat.value }}</div>
            </div>
            <div class="stat-icon" :style="{ background: stat.bgColor }">
              {{ stat.icon }}
            </div>
          </div>
          <div class="stat-footer">
            较昨日 <span :class="stat.trend > 0 ? 'up' : 'down'">
              {{ stat.trend > 0 ? '+' : '' }}{{ stat.trend }}%
            </span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 图表区域 -->
    <a-row :gutter="16" style="margin-top: 16px;">
      <a-col :span="16">
        <a-card title="本周收入趋势">
          <div class="chart-placeholder">
            <div class="bar-chart">
              <div class="bar" v-for="(v, i) in weekRevenue" :key="i"
                :style="{ height: v + '%', background: i === 4 ? '#C8E632' : '#2255CC' }">
                <span class="bar-label">{{ weekDays[i] }}</span>
              </div>
            </div>
          </div>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card title="收入构成">
          <div class="pie-placeholder">
            <div class="pie-item" v-for="item in revenueComposition" :key="item.name">
              <span class="pie-dot" :style="{ background: item.color }"></span>
              <span class="pie-name">{{ item.name }}</span>
              <span class="pie-value">{{ item.percent }}%</span>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 最近订单 -->
    <a-card title="最近订单" style="margin-top: 16px;">
      <a-table :data="recentOrders" :pagination="false" size="small">
        <template #columns>
          <a-table-column title="订单号" data-index="orderNo" />
          <a-table-column title="用户" data-index="user" />
          <a-table-column title="类型" data-index="type">
            <template #cell="{ record }">
              <a-tag :color="typeColor(record.type)">{{ record.type }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="金额" data-index="amount">
            <template #cell="{ record }">
              ¥{{ record.amount }}
            </template>
          </a-table-column>
          <a-table-column title="状态" data-index="status">
            <template #cell="{ record }">
              <a-tag :color="statusColor(record.status)">{{ record.status }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="时间" data-index="time" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const stats = ref([
  { title: '今日营收', value: '¥12,680', icon: '💰', color: '#2255CC', bgColor: '#EEF2FF', trend: 12 },
  { title: '今日订单', value: '48笔', icon: '📋', color: '#22C55E', bgColor: '#ECFDF5', trend: 8 },
  { title: '场地使用率', value: '78%', icon: '🎾', color: '#F59E0B', bgColor: '#FEF3C7', trend: -3 },
  { title: '活跃用户', value: '326人', icon: '👥', color: '#8B5CF6', bgColor: '#F5F3FF', trend: 15 },
])

const weekDays = ['一', '二', '三', '四', '五', '六', '日']
const weekRevenue = [45, 62, 55, 70, 90, 85, 78]

const revenueComposition = ref([
  { name: '场地预约', percent: 45, color: '#2255CC' },
  { name: '课程收入', percent: 30, color: '#22C55E' },
  { name: '商品销售', percent: 15, color: '#F59E0B' },
  { name: '活动报名', percent: 10, color: '#C8E632' },
])

const recentOrders = ref([
  { orderNo: '2026031900001', user: '张伟', type: '场地预约', amount: '280', status: '已支付', time: '10:30' },
  { orderNo: '2026031900002', user: '李娜', type: '课程购买', amount: '1,280', status: '已支付', time: '10:15' },
  { orderNo: '2026031900003', user: '王明', type: '活动报名', amount: '68', status: '已完成', time: '09:45' },
  { orderNo: '2026031900004', user: '刘洋', type: '商品购买', amount: '350', status: '待支付', time: '09:30' },
  { orderNo: '2026031900005', user: '陈静', type: '场地预约', amount: '140', status: '已取消', time: '09:10' },
])

function typeColor(type: string) {
  const map: Record<string, string> = { '场地预约': 'blue', '课程购买': 'green', '活动报名': 'orange', '商品购买': 'purple' }
  return map[type] || 'gray'
}

function statusColor(status: string) {
  const map: Record<string, string> = { '已支付': 'green', '已完成': 'blue', '待支付': 'orange', '已取消': 'red' }
  return map[status] || 'gray'
}
</script>

<style lang="scss" scoped>
.stat-card { border-radius: 8px; }
.stat-content { display: flex; justify-content: space-between; align-items: center; }
.stat-title { font-size: 14px; color: #999; }
.stat-value { font-size: 28px; font-weight: 700; margin-top: 8px; }
.stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 24px; }
.stat-footer { margin-top: 12px; font-size: 12px; color: #999; }
.up { color: #22C55E; }
.down { color: #EF4444; }

.chart-placeholder { height: 250px; display: flex; align-items: flex-end; }
.bar-chart { display: flex; align-items: flex-end; gap: 16px; width: 100%; height: 100%; padding: 20px 0; }
.bar {
  flex: 1; border-radius: 4px 4px 0 0; min-height: 10px;
  position: relative; transition: height 0.3s;
}
.bar-label {
  position: absolute; bottom: -24px; left: 50%; transform: translateX(-50%);
  font-size: 12px; color: #999;
}

.pie-placeholder { padding: 20px 0; }
.pie-item { display: flex; align-items: center; gap: 8px; margin-bottom: 16px; }
.pie-dot { width: 12px; height: 12px; border-radius: 50%; }
.pie-name { flex: 1; font-size: 14px; color: #333; }
.pie-value { font-size: 14px; font-weight: 600; color: #333; }
</style>
