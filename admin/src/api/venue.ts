import { get, post, put, del } from './request'

export const venueApi = {
  getVenues: (params?: any) => get('/admin/venues', params),
  createVenue: (data: any) => post('/admin/venues', data),
  updateVenue: (id: number, data: any) => put(`/admin/venues/${id}`, data),
  deleteVenue: (id: number) => del(`/admin/venues/${id}`),
  getTimeGrid: (params: { date: string }) => get('/admin/venues/time-grid', params),
  getTimeRules: (venueId: number) => get(`/admin/venues/${venueId}/time-rules`),
  createTimeRule: (venueId: number, data: any) => post(`/admin/venues/${venueId}/time-rules`, data),
}
