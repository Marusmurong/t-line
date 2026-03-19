-- Bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    venue_id BIGINT NOT NULL REFERENCES venues(id),
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    duration_hours DECIMAL(3,1) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    order_id BIGINT REFERENCES orders(id),
    cancel_reason VARCHAR(256) DEFAULT '',
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_venue_id ON bookings(venue_id);
CREATE INDEX idx_bookings_date ON bookings(date);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_order_id ON bookings(order_id);
CREATE INDEX idx_bookings_venue_date ON bookings(venue_id, date);

-- Booking waitlist table
CREATE TABLE IF NOT EXISTS booking_waitlist (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT REFERENCES bookings(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    venue_id BIGINT NOT NULL REFERENCES venues(id),
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    position INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting',
    notified_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_waitlist_venue_date ON booking_waitlist(venue_id, date);
CREATE INDEX idx_waitlist_user_id ON booking_waitlist(user_id);
CREATE INDEX idx_waitlist_status ON booking_waitlist(status);
