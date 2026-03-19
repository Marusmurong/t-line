import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/default.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/index.vue'), meta: { title: '总览' } },
      { path: 'venue/list', name: 'VenueList', component: () => import('@/views/venue/list.vue'), meta: { title: '场地列表' } },
      { path: 'venue/time-grid', name: 'VenueTimeGrid', component: () => import('@/views/venue/time-grid.vue'), meta: { title: '时段视图' } },
      { path: 'order/list', name: 'OrderList', component: () => import('@/views/order/list.vue'), meta: { title: '订单管理' } },
      { path: 'activity/list', name: 'ActivityList', component: () => import('@/views/activity/list.vue'), meta: { title: '活动管理' } },
      { path: 'product/list', name: 'ProductList', component: () => import('@/views/product/list.vue'), meta: { title: '商品管理' } },
      { path: 'academic/schedules', name: 'Schedules', component: () => import('@/views/academic/schedules.vue'), meta: { title: '课程管理' } },
      { path: 'academic/coaches', name: 'Coaches', component: () => import('@/views/academic/coaches.vue'), meta: { title: '教练管理' } },
      { path: 'academic/students', name: 'Students', component: () => import('@/views/academic/students.vue'), meta: { title: '学员管理' } },
      { path: 'stats', name: 'Stats', component: () => import('@/views/stats/index.vue'), meta: { title: '数据统计' } },
      { path: 'settings', name: 'Settings', component: () => import('@/views/settings/index.vue'), meta: { title: '系统设置' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// navigation guard
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.requiresAuth !== false && !token) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
