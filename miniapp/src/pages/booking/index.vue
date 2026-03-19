<template>
  <view class="booking">
    <!-- 日期选择器 -->
    <view class="date-picker">
      <view
        class="date-item"
        v-for="item in dateList"
        :key="item.dateStr"
        :class="{
          'date-item--active': item.dateStr === selectedDate,
          'date-item--weekend': item.isWeekend,
        }"
        @tap="onSelectDate(item.dateStr)"
      >
        <text class="date-weekday">{{ item.weekday }}</text>
        <text class="date-day">{{ item.day }}</text>
      </view>
    </view>

    <!-- 场地类型筛选 -->
    <view class="filter-tabs">
      <view
        class="filter-tab"
        v-for="tab in filterTabs"
        :key="tab.value"
        :class="{ 'filter-tab--active': activeFilter === tab.value }"
        @tap="onFilterChange(tab.value)"
      >
        <text>{{ tab.label }}</text>
      </view>
    </view>

    <!-- 时间轴头部 -->
    <view class="timeline-header">
      <view class="venue-label-placeholder" />
      <scroll-view scroll-x class="time-header-scroll" :scroll-left="timeScrollLeft">
        <view class="time-labels">
          <view class="time-label" v-for="h in timeHours" :key="h">
            <text>{{ h }}:00</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 场地列表 + 时间网格 -->
    <scroll-view scroll-y class="venue-list" :style="{ height: venueListHeight }">
      <view class="venue-row" v-for="venue in filteredVenues" :key="venue.id">
        <view class="venue-info">
          <text class="venue-name">{{ venue.name }}</text>
          <view
            class="venue-type-tag"
            :class="{
              'tag-indoor': venue.type === 'indoor',
              'tag-outdoor': venue.type === 'outdoor',
              'tag-practice': venue.type === 'practice',
            }"
          >
            <text>{{ venueTypeLabel(venue.type) }}</text>
          </view>
        </view>
        <scroll-view scroll-x class="time-grid-scroll" @scroll="onTimeGridScroll">
          <view class="time-grid">
            <view
              class="time-slot"
              v-for="slot in venue.slots"
              :key="slot.time"
              :class="slotClass(venue.id, slot)"
              @tap="onSlotTap(venue, slot)"
            >
              <text class="slot-text">{{ slotStatusText(venue.id, slot) }}</text>
            </view>
          </view>
        </scroll-view>
      </view>
    </scroll-view>

    <!-- 底部固定栏 -->
    <view class="bottom-bar" v-if="selectedSlots.length > 0">
      <view class="bottom-info">
        <text class="bottom-count">已选 {{ selectedSlots.length }} 个时段</text>
        <text class="bottom-price">合计 ¥{{ totalPrice }}</text>
      </view>
      <view class="bottom-btn" @tap="onConfirmBooking">
        <text class="bottom-btn-text">立即预订</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { venueApi } from '@/api/venue'

// ========== 类型 ==========
interface TimeSlot {
  time: string
  status: 'available' | 'booked' | 'maintenance'
  price: number
}

interface Venue {
  id: number
  name: string
  type: 'indoor' | 'outdoor' | 'practice'
  slots: TimeSlot[]
}

interface SelectedSlot {
  venueId: number
  venueName: string
  venueType: string
  time: string
  price: number
}

interface DateItem {
  dateStr: string
  day: string
  weekday: string
  isWeekend: boolean
}

// ========== 日期选择 ==========
function buildDateList(): DateItem[] {
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const list: DateItem[] = []
  for (let i = 0; i < 5; i++) {
    const d = new Date()
    d.setDate(d.getDate() + i)
    const month = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const dow = d.getDay()
    list.push({
      dateStr: `${d.getFullYear()}-${month}-${day}`,
      day: `${month}/${day}`,
      weekday: i === 0 ? '今天' : weekdays[dow],
      isWeekend: dow === 0 || dow === 6,
    })
  }
  return list
}

const dateList = ref<DateItem[]>(buildDateList())
const selectedDate = ref(dateList.value[0].dateStr)

function onSelectDate(date: string) {
  selectedDate.value = date
  loadAvailability()
}

// ========== 筛选 ==========
const filterTabs = [
  { label: '全部', value: 'all' },
  { label: '室内', value: 'indoor' },
  { label: '室外', value: 'outdoor' },
  { label: '练习', value: 'practice' },
]
const activeFilter = ref('all')

function onFilterChange(value: string) {
  activeFilter.value = value
}

// ========== 时间轴 ==========
const timeHours = Array.from({ length: 15 }, (_, i) => i + 8) // 8:00 ~ 22:00
const timeScrollLeft = ref(0)

function onTimeGridScroll(e: any) {
  timeScrollLeft.value = e.detail.scrollLeft
}

// ========== Mock 数据 ==========
function generateMockSlots(): TimeSlot[] {
  return timeHours.map((h) => {
    const rand = Math.random()
    let status: TimeSlot['status'] = 'available'
    if (rand > 0.7) status = 'booked'
    else if (rand > 0.6) status = 'maintenance'
    return {
      time: `${String(h).padStart(2, '0')}:00`,
      status,
      price: h >= 18 ? 180 : 140,
    }
  })
}

const venues = ref<Venue[]>([
  { id: 1, name: '1号场', type: 'indoor', slots: generateMockSlots() },
  { id: 2, name: '2号场', type: 'indoor', slots: generateMockSlots() },
  { id: 3, name: '3号场', type: 'outdoor', slots: generateMockSlots() },
  { id: 4, name: '4号场', type: 'outdoor', slots: generateMockSlots() },
  { id: 5, name: '练习场A', type: 'practice', slots: generateMockSlots() },
  { id: 6, name: '练习场B', type: 'practice', slots: generateMockSlots() },
])

const filteredVenues = computed(() => {
  if (activeFilter.value === 'all') return venues.value
  return venues.value.filter((v) => v.type === activeFilter.value)
})

// ========== 选择时段 ==========
const selectedSlots = ref<SelectedSlot[]>([])

function isSelected(venueId: number, time: string): boolean {
  return selectedSlots.value.some((s) => s.venueId === venueId && s.time === time)
}

function slotClass(venueId: number, slot: TimeSlot): string {
  if (isSelected(venueId, slot.time)) return 'slot-selected'
  if (slot.status === 'available') return 'slot-available'
  if (slot.status === 'booked') return 'slot-booked'
  return 'slot-maintenance'
}

function slotStatusText(venueId: number, slot: TimeSlot): string {
  if (isSelected(venueId, slot.time)) return '已选'
  if (slot.status === 'available') return '可约'
  if (slot.status === 'booked') return '已满'
  return '维护'
}

function onSlotTap(venue: Venue, slot: TimeSlot) {
  if (slot.status !== 'available') return

  const idx = selectedSlots.value.findIndex(
    (s) => s.venueId === venue.id && s.time === slot.time,
  )
  if (idx >= 0) {
    selectedSlots.value = [
      ...selectedSlots.value.slice(0, idx),
      ...selectedSlots.value.slice(idx + 1),
    ]
  } else {
    selectedSlots.value = [
      ...selectedSlots.value,
      {
        venueId: venue.id,
        venueName: venue.name,
        venueType: venue.type,
        time: slot.time,
        price: slot.price,
      },
    ]
  }
}

const totalPrice = computed(() =>
  selectedSlots.value.reduce((sum, s) => sum + s.price, 0),
)

// ========== 场地类型标签 ==========
function venueTypeLabel(type: string): string {
  const map: Record<string, string> = {
    indoor: '室内',
    outdoor: '室外',
    practice: '练习',
  }
  return map[type] || type
}

// ========== 列表高度计算 ==========
const venueListHeight = ref('calc(100vh - 400rpx)')

// ========== API 调用（保留结构，当前使用 mock） ==========
async function loadAvailability() {
  try {
    // TODO: 接入真实 API 后替换 mock 数据
    // const data = await venueApi.getAvailability(venueId, selectedDate.value)
    venues.value = venues.value.map((v) => ({
      ...v,
      slots: generateMockSlots(),
    }))
    selectedSlots.value = []
  } catch {
    uni.showToast({ title: '加载场地信息失败', icon: 'none' })
  }
}

// ========== 确认预订 ==========
function onConfirmBooking() {
  const params = encodeURIComponent(JSON.stringify({
    date: selectedDate.value,
    slots: selectedSlots.value,
    totalPrice: totalPrice.value,
  }))
  uni.navigateTo({
    url: `/pages/booking/confirm?data=${params}`,
  })
}

onMounted(() => {
  loadAvailability()
})
</script>

<style lang="scss" scoped>
.booking {
  background: $color-bg;
  min-height: 100vh;
  padding-bottom: 160rpx;
}

// ========== 日期选择器 ==========
.date-picker {
  display: flex;
  background: $color-bg-card;
  padding: $spacing-sm $spacing-base;
  gap: $spacing-sm;
}

.date-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: $spacing-sm 0;
  border-radius: $radius-base;
  gap: 4rpx;
  transition: all 0.2s;
}

.date-weekday {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.date-day {
  font-size: $font-size-base;
  font-weight: 600;
  color: $color-text;
}

.date-item--active {
  background: $brand-primary;
  .date-weekday,
  .date-day {
    color: #fff;
  }
}

.date-item--weekend:not(.date-item--active) {
  .date-weekday {
    color: $color-warning;
  }
}

// ========== 筛选标签 ==========
.filter-tabs {
  display: flex;
  padding: $spacing-sm $spacing-base;
  gap: $spacing-sm;
}

.filter-tab {
  padding: $spacing-xs $spacing-base;
  border-radius: $radius-round;
  background: $color-bg-card;
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.filter-tab--active {
  background: $brand-primary;
  color: #fff;
}

// ========== 时间轴头部 ==========
.timeline-header {
  display: flex;
  background: $color-bg-card;
  border-bottom: 1rpx solid $color-border;
}

.venue-label-placeholder {
  width: 160rpx;
  flex-shrink: 0;
}

.time-header-scroll {
  flex: 1;
  white-space: nowrap;
}

.time-labels {
  display: flex;
}

.time-label {
  width: 100rpx;
  flex-shrink: 0;
  text-align: center;
  font-size: $font-size-xs;
  color: $color-text-placeholder;
  padding: $spacing-xs 0;
}

// ========== 场地行 ==========
.venue-list {
  background: $color-bg-card;
}

.venue-row {
  display: flex;
  border-bottom: 1rpx solid $color-border;
}

.venue-info {
  width: 160rpx;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-sm $spacing-xs;
  gap: 8rpx;
}

.venue-name {
  font-size: $font-size-sm;
  font-weight: 600;
  color: $color-text;
}

.venue-type-tag {
  font-size: $font-size-xs;
  padding: 2rpx 12rpx;
  border-radius: $radius-sm;
}

.tag-indoor {
  background: #EEF2FF;
  color: $brand-primary;
}

.tag-outdoor {
  background: #ECFDF5;
  color: $color-success;
}

.tag-practice {
  background: #FEF3C7;
  color: $color-warning;
}

// ========== 时间网格 ==========
.time-grid-scroll {
  flex: 1;
  white-space: nowrap;
}

.time-grid {
  display: flex;
}

.time-slot {
  width: 100rpx;
  height: 80rpx;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-right: 1rpx solid $color-border;
}

.slot-text {
  font-size: $font-size-xs;
}

.slot-available {
  background: #ECFDF5;
  .slot-text { color: $color-success; }
}

.slot-booked {
  background: #FEF2F2;
  .slot-text { color: $color-error; }
}

.slot-maintenance {
  background: #FEF3C7;
  .slot-text { color: $color-warning; }
}

.slot-selected {
  background: $brand-primary;
  .slot-text { color: #fff; }
}

// ========== 底部栏 ==========
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-sm $spacing-base;
  padding-bottom: calc(#{$spacing-sm} + env(safe-area-inset-bottom));
  background: $color-bg-card;
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.06);
}

.bottom-info {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.bottom-count {
  font-size: $font-size-sm;
  color: $color-text-secondary;
}

.bottom-price {
  font-size: $font-size-lg;
  font-weight: 700;
  color: $color-error;
}

.bottom-btn {
  background: $brand-primary;
  padding: $spacing-sm $spacing-lg;
  border-radius: $radius-round;
}

.bottom-btn-text {
  color: #fff;
  font-size: $font-size-base;
  font-weight: 600;
}
</style>
