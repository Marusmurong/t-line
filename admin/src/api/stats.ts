import { get } from './request'

export const statsApi = {
  getDashboard: () => get('/admin/dashboard'),

  getRevenueStats: (params: { start_date: string; end_date: string }) =>
    get('/admin/stats/revenue', params),

  getVenueUsageStats: (params: { start_date: string; end_date: string }) =>
    get('/admin/stats/venue-usage', params),

  getUserStats: (params: { start_date: string; end_date: string }) =>
    get('/admin/stats/users', params),
}
