import { get, post } from './request'

export const activityApi = {
  getActivities: (params?: any) =>
    get('/activities', params),

  getActivityDetail: (id: number) =>
    get(`/activities/${id}`),

  register: (id: number) =>
    post(`/activities/${id}/register`),

  cancelRegistration: (id: number) =>
    post(`/activities/${id}/cancel-registration`),
}
