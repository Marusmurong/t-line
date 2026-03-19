<template>
  <div class="time-grid-page">
    <!-- 顶部筛选 -->
    <div class="page-header">
      <h2 class="page-title">时段视图</h2>
      <div class="header-filters">
        <a-date-picker
          v-model="selectedDate"
          style="width: 180px"
          placeholder="选择日期"
        />
        <a-select
          v-model="selectedVenue"
          placeholder="场地筛选"
          allow-clear
          style="width: 160px"
        >
          <a-option value="">全部场地</a-option>
          <a-option v-for="v in venueNames" :key="v" :value="v">{{ v }}</a-option>
        </a-select>
      </div>
    </div>

    <!-- 图例 -->
    <div class="legend">
      <span class="legend-item"><span class="legend-dot dot-free"></span>空闲</span>
      <span class="legend-item"><span class="legend-dot dot-booked"></span>已预约</span>
      <span class="legend-item"><span class="legend-dot dot-expired"></span>已过期</span>
      <span class="legend-item"><span class="legend-dot dot-maintenance"></span>维护</span>
    </div>

    <div class="grid-layout">
      <!-- 时间网格主体 -->
      <a-card class="grid-card">
        <div class="time-grid">
          <!-- 表头：时间轴 -->
          <div class="grid-header">
            <div class="grid-label-cell">场地</div>
            <div class="grid-time-cell" v-for="h in timeSlots" :key="h">
              {{ h }}
            </div>
          </div>
          <!-- 每个场地一行 -->
          <div
            class="grid-row"
            v-for="venue in filteredGridData"
            :key="venue.name"
          >
            <div class="grid-label-cell">{{ venue.name }}</div>
            <div
              v-for="slot in venue.slots"
              :key="slot.time"
              :class="['grid-slot', `slot-${slot.status}`]"
              @click="handleSlotClick(venue.name, slot)"
            >
              <a-tooltip v-if="slot.status === 'booked'" :content="`${slot.user} ${slot.time}`">
                <div class="slot-inner"></div>
              </a-tooltip>
            </div>
          </div>
        </div>
      </a-card>

      <!-- 右侧详情面板 -->
      <a-card v-if="selectedSlot" class="detail-panel">
        <template #title>预约详情</template>
        <a-descriptions :column="1" size="small" bordered>
          <a-descriptions-item label="场地">{{ selectedSlot.venueName }}</a-descriptions-item>
          <a-descriptions-item label="时段">{{ selectedSlot.time }}</a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag :color="slotStatusColor[selectedSlot.status]">
              {{ slotStatusLabel[selectedSlot.status] }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item v-if="selectedSlot.user" label="预约用户">
            {{ selectedSlot.user }}
          </a-descriptions-item>
          <a-descriptions-item v-if="selectedSlot.phone" label="联系电话">
            {{ selectedSlot.phone }}
          </a-descriptions-item>
          <a-descriptions-item v-if="selectedSlot.orderNo" label="订单号">
            {{ selectedSlot.orderNo }}
          </a-descriptions-item>
        </a-descriptions>
        <a-button
          v-if="selectedSlot.status === 'booked'"
          type="outline"
          status="danger"
          size="small"
          style="margin-top: 12px; width: 100%"
        >
          取消预约
        </a-button>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface TimeSlot {
  time: string
  status: 'free' | 'booked' | 'expired' | 'maintenance'
  user?: string
  phone?: string
  orderNo?: string
}

interface SelectedSlotInfo extends TimeSlot {
  venueName: string
}

interface VenueGridRow {
  name: string
  slots: TimeSlot[]
}

const selectedDate = ref('2026-03-20')
const selectedVenue = ref('')

const venueNames = ['网球场A', '网球场B', '室内馆1号场', '室内馆2号场', '训练场C', 'VIP球场']

const timeSlots = Array.from({ length: 14 }, (_, i) => `${i + 8}:00`)

const slotStatusColor: Record<string, string> = {
  free: 'green',
  booked: 'blue',
  expired: 'red',
  maintenance: 'orangered',
}

const slotStatusLabel: Record<string, string> = {
  free: '空闲',
  booked: '已预约',
  expired: '已过期',
  maintenance: '维护',
}

function createSlots(): TimeSlot[] {
  return timeSlots.map((time) => ({ time, status: 'free' as const }))
}

function withBooking(
  slots: TimeSlot[],
  timeIndex: number,
  user: string,
  phone: string,
  orderNo: string,
  status: TimeSlot['status'] = 'booked',
): TimeSlot[] {
  return slots.map((s, i) =>
    i === timeIndex ? { ...s, status, user, phone, orderNo } : s,
  )
}

const gridData = ref<VenueGridRow[]>([
  {
    name: '网球场A',
    slots: withBooking(
      withBooking(
        withBooking(createSlots(), 0, '张伟', '138****1234', 'ORD20260320001', 'expired'),
        2, '李娜', '139****5678', 'ORD20260320002',
      ),
      6, '王明', '137****9012', 'ORD20260320003',
    ),
  },
  {
    name: '网球场B',
    slots: withBooking(
      withBooking(createSlots(), 4, '刘洋', '136****3456', 'ORD20260320004'),
      9, '陈静', '135****7890', 'ORD20260320005',
    ),
  },
  {
    name: '室内馆1号场',
    slots: withBooking(
      withBooking(
        withBooking(createSlots(), 1, '赵磊', '133****2345', 'ORD20260320006'),
        3, '孙芳', '132****6789', 'ORD20260320007',
      ),
      7, '周杰', '131****0123', 'ORD20260320008',
    ),
  },
  {
    name: '室内馆2号场',
    slots: createSlots().map((s) => ({ ...s, status: 'maintenance' as const })),
  },
  {
    name: '训练场C',
    slots: withBooking(createSlots(), 5, '吴强', '130****4567', 'ORD20260320009'),
  },
  {
    name: 'VIP球场',
    slots: withBooking(
      withBooking(createSlots(), 2, '郑华', '129****8901', 'ORD20260320010'),
      10, '黄丽', '128****2345', 'ORD20260320011',
    ),
  },
])

const filteredGridData = computed(() => {
  if (!selectedVenue.value) return gridData.value
  return gridData.value.filter((v) => v.name === selectedVenue.value)
})

const selectedSlot = ref<SelectedSlotInfo | null>(null)

function handleSlotClick(venueName: string, slot: TimeSlot) {
  selectedSlot.value = { ...slot, venueName }
}
</script>

<style lang="scss" scoped>
.time-grid-page {
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
    gap: 12px;
  }

  .legend {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
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

  .dot-free { background: #22c55e; }
  .dot-booked { background: #2255cc; }
  .dot-expired { background: #ef4444; }
  .dot-maintenance { background: #f59e0b; }

  .grid-layout {
    display: flex;
    gap: 16px;
  }

  .grid-card {
    flex: 1;
    overflow-x: auto;
  }

  .detail-panel {
    width: 280px;
    flex-shrink: 0;
  }

  .time-grid {
    min-width: 900px;
  }

  .grid-header,
  .grid-row {
    display: grid;
    grid-template-columns: 120px repeat(14, 1fr);
    gap: 2px;
  }

  .grid-header {
    margin-bottom: 4px;
  }

  .grid-label-cell {
    font-size: 13px;
    font-weight: 500;
    color: #333;
    padding: 8px 4px;
    display: flex;
    align-items: center;
  }

  .grid-time-cell {
    font-size: 11px;
    color: #999;
    text-align: center;
    padding: 4px 0;
  }

  .grid-slot {
    height: 40px;
    border-radius: 4px;
    cursor: pointer;
    transition: opacity 0.2s;

    &:hover {
      opacity: 0.8;
    }

    .slot-inner {
      width: 100%;
      height: 100%;
    }
  }

  .slot-free { background: #dcfce7; }
  .slot-booked { background: #2255cc; }
  .slot-expired { background: #ef4444; }
  .slot-maintenance { background: #f59e0b; }
}
</style>
