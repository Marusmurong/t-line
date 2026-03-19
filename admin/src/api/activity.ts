import { get, post, put, del } from './request'

export const activityApi = {
  getActivities: (params?: any) => get('/admin/activities', params),
  createActivity: (data: any) => post('/admin/activities', data),
  updateActivity: (id: number, data: any) => put(`/admin/activities/${id}`, data),
  deleteActivity: (id: number) => del(`/admin/activities/${id}`),
}
