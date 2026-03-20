import { get, post, put, del } from './request'

export const academicApi = {
  getCoaches: (params?: any) => get('/admin/coaches', params),
  createCoach: (data: any) => post('/admin/coaches', data),
  updateCoach: (id: number, data: any) => put(`/admin/coaches/${id}`, data),
  getCoachPerformance: (id: number) => get(`/admin/coaches/${id}/performance`),
  getSchedules: (params?: any) => get('/admin/schedules', params),
  createSchedule: (data: any) => post('/admin/schedules', data),
  checkConflict: (data: any) => post('/admin/schedules/conflict-check', data),
  getStudents: (params?: any) => get('/admin/students', params),
  getStudentRecords: (id: number) => get(`/admin/students/${id}/records`),
}
