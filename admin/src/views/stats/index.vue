<template>
  <div class="stats-page">
    <!-- 顶部时间范围选择器 -->
    <div class="page-header">
      <h2 class="page-title">数据统计</h2>
      <div class="header-filters">
        <a-radio-group v-model="timeRange" type="button" size="small">
          <a-radio value="week">本周</a-radio>
          <a-radio value="month">本月</a-radio>
          <a-radio value="quarter">本季</a-radio>
          <a-radio value="custom">自定义</a-radio>
        </a-radio-group>
        <a-range-picker
          v-if="timeRange === 'custom'"
          v-model="customRange"
          style="width: 260px"
        />
      </div>
    </div>

    <!-- KPI 卡片 -->
    <a-row :gutter="16">
      <a-col :span="6" v-for="kpi in kpiCards" :key="kpi.title">
        <a-card class="kpi-card">
          <div class="kpi-content">
            <div class="kpi-info">
              <div class="kpi-title">{{ kpi.title }}</div>
              <div class="kpi-value">{{ kpi.value }}</div>
            </div>
            <div class="kpi-icon" :style="{ background: kpi.bgColor, color: kpi.color }">
              <span class="kpi-icon-text">{{ kpi.icon }}</span>
            </div>
          </div>
          <div class="kpi-footer">
            同比
            <span :class="kpi.trend >= 0 ? 'trend-up' : 'trend-down'">
              {{ kpi.trend >= 0 ? '↑' : '↓' }}{{ Math.abs(kpi.trend) }}%
            </span>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 上半区：收入趋势 + 场地使用率热力图 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <!-- 月度收入趋势（柱状图） -->
      <a-col :span="14">
        <a-card title="月度收入趋势">
          <div class="bar-chart-wrap">
            <div class="bar-chart">
              <div
                v-for="(bar, idx) in revenueData"
                :key="idx"
                class="bar-col"
              >
                <div
                  class="bar"
                  :style="{
                    height: barHeight(bar.value) + '%',
                    background: bar.value === peakRevenue ? '#C8E632' : '#2255CC',
                  }"
                >
                  <span class="bar-tooltip">¥{{ bar.value.toLocaleString() }}</span>
                </div>
                <span v-if="idx % 5 === 0" class="bar-date">{{ bar.date }}</span>
                <span v-else class="bar-date bar-date-hidden">&nbsp;</span>
              </div>
            </div>
            <div class="bar-summary">
              日均收入：<strong>¥{{ dailyAvgRevenue.toLocaleString() }}</strong>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 场地使用率热力图 -->
      <a-col :span="10">
        <a-card title="场地使用率热力图">
          <div class="heatmap-wrap">
            <!-- 表头：星期 -->
            <div class="heatmap-grid heatmap-header">
              <div class="heatmap-label"></div>
              <div v-for="day in weekDays" :key="day" class="heatmap-head-cell">
                {{ day }}
              </div>
            </div>
            <!-- 行：时段 -->
            <div
              v-for="period in heatmapData"
              :key="period.label"
              class="heatmap-grid"
            >
              <div class="heatmap-label">{{ period.label }}</div>
              <div
                v-for="(val, dIdx) in period.values"
                :key="dIdx"
                class="heatmap-cell"
                :style="{ background: heatColor(val) }"
              >
                {{ val }}%
              </div>
            </div>
            <!-- 色阶图例 -->
            <div class="heatmap-legend">
              <span class="legend-text">低</span>
              <div class="legend-bar">
                <span
                  v-for="(c, i) in heatColorScale"
                  :key="i"
                  class="legend-swatch"
                  :style="{ background: c }"
                ></span>
              </div>
              <span class="legend-text">高</span>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 下半区：用户增长趋势 + 会员等级分布 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <!-- 用户增长趋势（面积图） -->
      <a-col :span="12">
        <a-card title="用户增长趋势">
          <div class="area-chart-wrap">
            <svg
              class="area-chart-svg"
              :viewBox="`0 0 ${areaWidth} ${areaHeight}`"
              preserveAspectRatio="none"
            >
              <defs>
                <linearGradient id="areaGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#2255CC" stop-opacity="0.3" />
                  <stop offset="100%" stop-color="#2255CC" stop-opacity="0.02" />
                </linearGradient>
              </defs>
              <path :d="areaPath" fill="url(#areaGrad)" />
              <path :d="linePath" fill="none" stroke="#2255CC" stroke-width="2" />
              <!-- 数据点 -->
              <circle
                v-for="(pt, i) in areaPoints"
                :key="i"
                :cx="pt.x"
                :cy="pt.y"
                r="3"
                fill="#2255CC"
              />
            </svg>
            <!-- X轴日期标签 -->
            <div class="area-x-labels">
              <span
                v-for="(d, i) in userGrowthData"
                :key="i"
                class="area-x-label"
                :style="{ visibility: i % 5 === 0 ? 'visible' : 'hidden' }"
              >
                {{ d.date }}
              </span>
            </div>
            <div class="area-summary">
              总计：<strong>{{ userGrowthTotal }}人</strong>
              &nbsp;&nbsp;日均：<strong>{{ userGrowthDailyAvg }}人</strong>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 会员等级分布（环形图） -->
      <a-col :span="12">
        <a-card title="会员等级分布">
          <div class="donut-wrap">
            <div class="donut-chart" :style="{ background: donutGradient }">
              <div class="donut-inner">
                <div class="donut-total-label">总会员</div>
                <div class="donut-total-value">{{ memberTotal }}</div>
              </div>
            </div>
            <div class="donut-legend">
              <div
                v-for="seg in memberSegments"
                :key="seg.label"
                class="donut-legend-item"
              >
                <span class="donut-dot" :style="{ background: seg.color }"></span>
                <span class="donut-label">{{ seg.label }}</span>
                <span class="donut-percent">{{ seg.percent }}%</span>
                <span class="donut-count">{{ seg.count }}人</span>
              </div>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

// ====== 时间范围 ======

const timeRange = ref<string>('month')
const customRange = ref<string[]>([])

// ====== KPI 卡片 ======

interface KpiItem {
  readonly title: string
  readonly value: string
  readonly icon: string
  readonly color: string
  readonly bgColor: string
  readonly trend: number
}

const kpiCards = ref<readonly KpiItem[]>([
  { title: '总收入', value: '¥89,500', icon: '💰', color: '#2255CC', bgColor: '#EEF2FF', trend: 12 },
  { title: '总订单', value: '368笔', icon: '📋', color: '#22C55E', bgColor: '#ECFDF5', trend: 8 },
  { title: '新增用户', value: '52人', icon: '👥', color: '#8B5CF6', bgColor: '#F5F3FF', trend: 23 },
  { title: '会员转化率', value: '23%', icon: '🏅', color: '#F59E0B', bgColor: '#FEF3C7', trend: -2 },
])

// ====== 月度收入趋势 ======

interface RevenueBar {
  readonly date: string
  readonly value: number
}

function generateRevenueData(): readonly RevenueBar[] {
  const base = [
    2800, 3200, 2900, 3500, 4100, 3800, 3100,
    2600, 3400, 3900, 4200, 3700, 3300, 2800,
    3600, 4500, 4800, 4100, 3500, 3000, 2700,
    3100, 3800, 4300, 4600, 5200, 4800, 4200,
    3600, 3100,
  ]
  return base.map((v, i) => ({
    date: `${i + 1}日`,
    value: v,
  }))
}

const revenueData = ref<readonly RevenueBar[]>(generateRevenueData())

const peakRevenue = computed(() =>
  Math.max(...revenueData.value.map((d) => d.value)),
)

const dailyAvgRevenue = computed(() => {
  const total = revenueData.value.reduce((sum, d) => sum + d.value, 0)
  return Math.round(total / revenueData.value.length)
})

function barHeight(value: number): number {
  const max = peakRevenue.value
  if (max === 0) return 0
  return Math.round((value / max) * 100)
}

// ====== 场地使用率热力图 ======

const weekDays = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

interface HeatmapRow {
  readonly label: string
  readonly values: readonly number[]
}

const heatmapData = ref<readonly HeatmapRow[]>([
  { label: '8:00-10:00', values: [35, 42, 38, 45, 50, 78, 82] },
  { label: '10:00-12:00', values: [55, 60, 58, 62, 68, 88, 92] },
  { label: '12:00-14:00', values: [30, 28, 32, 35, 40, 55, 60] },
  { label: '14:00-16:00', values: [48, 52, 50, 55, 65, 85, 90] },
  { label: '16:00-18:00', values: [62, 68, 65, 70, 75, 92, 95] },
  { label: '18:00-20:00', values: [75, 80, 78, 82, 88, 95, 98] },
  { label: '20:00-22:00', values: [68, 72, 70, 75, 80, 88, 85] },
])

const heatColorScale = ['#EEF2FF', '#C7D2FE', '#93A8F4', '#6380E8', '#3B5BD9', '#2255CC']

function heatColor(value: number): string {
  if (value < 30) return heatColorScale[0]
  if (value < 45) return heatColorScale[1]
  if (value < 60) return heatColorScale[2]
  if (value < 75) return heatColorScale[3]
  if (value < 90) return heatColorScale[4]
  return heatColorScale[5]
}

// ====== 用户增长趋势（面积图） ======

interface UserGrowthPoint {
  readonly date: string
  readonly value: number
}

function generateUserGrowthData(): readonly UserGrowthPoint[] {
  const values = [
    1, 3, 2, 1, 2, 4, 3, 2, 1, 0,
    2, 3, 1, 2, 1, 3, 2, 4, 2, 1,
    0, 2, 1, 3, 2, 1, 2, 3, 1, 2,
  ]
  return values.map((v, i) => ({
    date: `${i + 1}日`,
    value: v,
  }))
}

const userGrowthData = ref<readonly UserGrowthPoint[]>(generateUserGrowthData())

const userGrowthTotal = computed(() =>
  userGrowthData.value.reduce((sum, d) => sum + d.value, 0),
)

const userGrowthDailyAvg = computed(() =>
  Math.round((userGrowthTotal.value / userGrowthData.value.length) * 10) / 10,
)

const areaWidth = 600
const areaHeight = 200
const areaPadding = 20

interface ChartPoint {
  readonly x: number
  readonly y: number
}

const areaPoints = computed<readonly ChartPoint[]>(() => {
  const data = userGrowthData.value
  const maxVal = Math.max(...data.map((d) => d.value), 1)
  const usableWidth = areaWidth - areaPadding * 2
  const usableHeight = areaHeight - areaPadding * 2

  return data.map((d, i) => ({
    x: areaPadding + (i / (data.length - 1)) * usableWidth,
    y: areaPadding + usableHeight - (d.value / maxVal) * usableHeight,
  }))
})

const linePath = computed(() => {
  const pts = areaPoints.value
  if (pts.length === 0) return ''
  return pts.map((p, i) => `${i === 0 ? 'M' : 'L'}${p.x},${p.y}`).join(' ')
})

const areaPath = computed(() => {
  const pts = areaPoints.value
  if (pts.length === 0) return ''
  const bottomY = areaHeight - areaPadding
  const firstX = pts[0].x
  const lastX = pts[pts.length - 1].x
  return (
    `M${firstX},${bottomY} ` +
    pts.map((p) => `L${p.x},${p.y}`).join(' ') +
    ` L${lastX},${bottomY} Z`
  )
})

// ====== 会员等级分布（环形图） ======

interface MemberSegment {
  readonly label: string
  readonly percent: number
  readonly count: number
  readonly color: string
}

const memberSegments = ref<readonly MemberSegment[]>([
  { label: '普通会员', percent: 60, count: 1200, color: '#9CA3AF' },
  { label: '银卡会员', percent: 20, count: 400, color: '#C0C0C0' },
  { label: '金卡会员', percent: 15, count: 300, color: '#F59E0B' },
  { label: '钻石会员', percent: 5, count: 100, color: '#2255CC' },
])

const memberTotal = computed(() =>
  memberSegments.value.reduce((sum, s) => sum + s.count, 0),
)

const donutGradient = computed(() => {
  const segments = memberSegments.value
  const parts: string[] = []
  let cumulative = 0
  for (const seg of segments) {
    const start = cumulative
    const end = cumulative + seg.percent
    parts.push(`${seg.color} ${start}% ${end}%`)
    cumulative = end
  }
  return `conic-gradient(${parts.join(', ')})`
})
</script>

<style lang="scss" scoped>
.stats-page {
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

  .header-filters {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

// ====== KPI 卡片 ======

.kpi-card {
  border-radius: 8px;
}

.kpi-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.kpi-title {
  font-size: 14px;
  color: #999;
}

.kpi-value {
  font-size: 28px;
  font-weight: 700;
  margin-top: 8px;
  color: #1d2129;
}

.kpi-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.kpi-icon-text {
  font-size: 24px;
}

.kpi-footer {
  margin-top: 12px;
  font-size: 12px;
  color: #999;
}

.trend-up {
  color: #22c55e;
  font-weight: 600;
}

.trend-down {
  color: #ef4444;
  font-weight: 600;
}

// ====== 柱状图 ======

.bar-chart-wrap {
  padding: 8px 0;
}

.bar-chart {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 220px;
  padding-bottom: 28px;
}

.bar-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  height: 100%;
  position: relative;
}

.bar {
  width: 100%;
  min-height: 4px;
  border-radius: 3px 3px 0 0;
  transition: height 0.3s ease;
  position: relative;
  cursor: pointer;

  &:hover .bar-tooltip {
    opacity: 1;
    transform: translateX(-50%) translateY(-4px);
  }
}

.bar-tooltip {
  position: absolute;
  top: -28px;
  left: 50%;
  transform: translateX(-50%) translateY(0);
  font-size: 11px;
  color: #fff;
  background: rgba(0, 0, 0, 0.75);
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
  opacity: 0;
  transition: opacity 0.2s, transform 0.2s;
  pointer-events: none;
}

.bar-date {
  position: absolute;
  bottom: -24px;
  font-size: 11px;
  color: #999;
  white-space: nowrap;
}

.bar-date-hidden {
  visibility: hidden;
}

.bar-summary {
  margin-top: 12px;
  font-size: 13px;
  color: #666;
  text-align: center;

  strong {
    color: #2255cc;
  }
}

// ====== 热力图 ======

.heatmap-wrap {
  padding: 4px 0;
}

.heatmap-grid {
  display: grid;
  grid-template-columns: 100px repeat(7, 1fr);
  gap: 3px;
  margin-bottom: 3px;
}

.heatmap-header {
  margin-bottom: 6px;
}

.heatmap-label {
  font-size: 12px;
  color: #666;
  display: flex;
  align-items: center;
  padding-right: 8px;
  white-space: nowrap;
}

.heatmap-head-cell {
  font-size: 12px;
  color: #999;
  text-align: center;
  padding: 4px 0;
}

.heatmap-cell {
  text-align: center;
  font-size: 11px;
  font-weight: 500;
  color: #fff;
  padding: 10px 0;
  border-radius: 4px;
  transition: transform 0.15s;

  &:hover {
    transform: scale(1.08);
    z-index: 1;
  }
}

.heatmap-legend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 14px;
}

.legend-text {
  font-size: 11px;
  color: #999;
}

.legend-bar {
  display: flex;
  gap: 2px;
}

.legend-swatch {
  width: 28px;
  height: 14px;
  border-radius: 2px;
}

// ====== 面积图 ======

.area-chart-wrap {
  padding: 8px 0;
}

.area-chart-svg {
  width: 100%;
  height: 200px;
}

.area-x-labels {
  display: flex;
  justify-content: space-between;
  padding: 0 20px;
  margin-top: 4px;
}

.area-x-label {
  font-size: 11px;
  color: #999;
  flex: 1;
  text-align: center;
}

.area-summary {
  margin-top: 12px;
  font-size: 13px;
  color: #666;
  text-align: center;

  strong {
    color: #2255cc;
  }
}

// ====== 环形图 ======

.donut-wrap {
  display: flex;
  align-items: center;
  gap: 32px;
  padding: 16px 0;
}

.donut-chart {
  width: 180px;
  height: 180px;
  border-radius: 50%;
  position: relative;
  flex-shrink: 0;
}

.donut-inner {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100px;
  height: 100px;
  background: #fff;
  border-radius: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.donut-total-label {
  font-size: 12px;
  color: #999;
}

.donut-total-value {
  font-size: 22px;
  font-weight: 700;
  color: #1d2129;
  margin-top: 2px;
}

.donut-legend {
  flex: 1;
}

.donut-legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;

  &:last-child {
    margin-bottom: 0;
  }
}

.donut-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.donut-label {
  font-size: 14px;
  color: #333;
  min-width: 72px;
}

.donut-percent {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  min-width: 40px;
  text-align: right;
}

.donut-count {
  font-size: 13px;
  color: #999;
  min-width: 50px;
  text-align: right;
}
</style>
