import { get, post } from './request'

export const academicApi = {
  getCoaches: (params?: any) => get('/coaches', params),
  getCoachDetail: (id: number) => get(`/coaches/${id}`),
  getMyCourseRecords: (params?: any) => get('/my-courses/records', params),
  rateRecord: (id: number, data: any) => post(`/records/${id}/rating`, data),
}
