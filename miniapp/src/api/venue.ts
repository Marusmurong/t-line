import { get } from './request'

export const venueApi = {
  getVenues: (params?: { type?: string }) =>
    get('/venues', params),

  getVenueDetail: (id: number) =>
    get(`/venues/${id}`),

  getAvailability: (id: number, date: string) =>
    get(`/venues/${id}/availability`, { date }),
}
