import { get, post } from './request'

export const bookingApi = {
  createBooking: (data: {
    venue_id: number
    date: string
    time_slots: string[]
  }) => post('/bookings', data),

  getBookings: (params?: { status?: string; page?: number }) =>
    get('/bookings', params),

  getBookingDetail: (id: number) =>
    get(`/bookings/${id}`),

  cancelBooking: (id: number) =>
    post(`/bookings/${id}/cancel`),

  joinWaitlist: (data: {
    venue_id: number
    date: string
    time_slot: string
  }) => post('/bookings/waitlist', data),
}
