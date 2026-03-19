<template>
  <a-layout class="layout">
    <!-- 侧边栏 -->
    <a-layout-sider :width="240" class="sider" collapsible :collapsed="collapsed" @collapse="collapsed = !collapsed">
      <div class="logo">
        <span class="logo-text" v-if="!collapsed">T-LINE</span>
        <span class="logo-text" v-else>T</span>
      </div>

      <a-menu
        :selected-keys="[currentRoute]"
        @menu-item-click="handleMenuClick"
        auto-open-selected
      >
        <a-menu-item key="dashboard">
          <template #icon><icon-dashboard /></template>
          总览
        </a-menu-item>

        <a-sub-menu key="venue">
          <template #icon><icon-location /></template>
          <template #title>场地管理</template>
          <a-menu-item key="venue/list">场地列表</a-menu-item>
          <a-menu-item key="venue/time-grid">时段视图</a-menu-item>
        </a-sub-menu>

        <a-menu-item key="activity/list">
          <template #icon><icon-trophy /></template>
          活动管理
        </a-menu-item>

        <a-menu-item key="product/list">
          <template #icon><icon-apps /></template>
          商品管理
        </a-menu-item>

        <a-menu-item key="order/list">
          <template #icon><icon-file /></template>
          订单管理
        </a-menu-item>

        <a-sub-menu key="academic">
          <template #icon><icon-book /></template>
          <template #title>教务管理</template>
          <a-menu-item key="academic/schedules">课程管理</a-menu-item>
          <a-menu-item key="academic/coaches">教练管理</a-menu-item>
          <a-menu-item key="academic/students">学员管理</a-menu-item>
        </a-sub-menu>

        <a-menu-item key="stats">
          <template #icon><icon-bar-chart /></template>
          数据统计
        </a-menu-item>

        <a-menu-item key="settings">
          <template #icon><icon-settings /></template>
          系统设置
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <!-- 右侧内容 -->
    <a-layout>
      <!-- 顶部栏 -->
      <a-layout-header class="header">
        <a-breadcrumb>
          <a-breadcrumb-item>T-Line 管理后台</a-breadcrumb-item>
          <a-breadcrumb-item>{{ currentTitle }}</a-breadcrumb-item>
        </a-breadcrumb>
        <div class="header-right">
          <a-badge :count="3" dot>
            <a-button type="text" shape="circle">
              <icon-notification />
            </a-button>
          </a-badge>
          <a-dropdown>
            <a-button type="text">
              <a-avatar :size="28" style="background-color: #2255CC">管</a-avatar>
              <span style="margin-left: 8px">管理员</span>
            </a-button>
            <template #content>
              <a-doption @click="handleLogout">退出登录</a-doption>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- 内容区 -->
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  IconDashboard, IconLocation, IconTrophy, IconApps,
  IconFile, IconBook, IconBarChart, IconSettings, IconNotification,
} from '@arco-design/web-vue/es/icon'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)

const currentRoute = computed(() => {
  const path = route.path.replace(/^\//, '')
  return path || 'dashboard'
})

const currentTitle = computed(() => {
  return (route.meta?.title as string) || '总览'
})

function handleMenuClick(key: string) {
  router.push(`/${key}`)
}

function handleLogout() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  router.push('/login')
}
</script>

<style lang="scss" scoped>
.layout { min-height: 100vh; }

.sider {
  :deep(.arco-layout-sider-children) {
    display: flex; flex-direction: column;
  }
}

.logo {
  height: 64px; display: flex; align-items: center; justify-content: center;
  background: #1a1f36; color: #fff; font-weight: 800; font-size: 22px;
  letter-spacing: 2px;
}

.header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 24px; background: #fff;
  border-bottom: 1px solid #e5e7eb;
}

.header-right { display: flex; align-items: center; gap: 8px; }

.content {
  padding: 24px; background: #f7f8fa; min-height: calc(100vh - 64px);
}
</style>
