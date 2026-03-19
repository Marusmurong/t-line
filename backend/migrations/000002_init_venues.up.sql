-- Venues table
CREATE TABLE IF NOT EXISTS venues (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(20) NOT NULL,
    description TEXT DEFAULT '',
    cover_image VARCHAR(512) DEFAULT '',
    facilities JSONB DEFAULT '{}',
    status SMALLINT DEFAULT 1,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_venues_status ON venues(status);
CREATE INDEX idx_venues_sort_order ON venues(sort_order);

-- Venue time slot rules table
CREATE TABLE IF NOT EXISTS venue_time_slot_rules (
    id BIGSERIAL PRIMARY KEY,
    venue_id BIGINT NOT NULL REFERENCES venues(id),
    day_type VARCHAR(20) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    member_discount JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_venue_time_rules_venue_id ON venue_time_slot_rules(venue_id);
CREATE INDEX idx_venue_time_rules_day_type ON venue_time_slot_rules(day_type);

-- Venue blocked times table
CREATE TABLE IF NOT EXISTS venue_blocked_times (
    id BIGSERIAL PRIMARY KEY,
    venue_id BIGINT NOT NULL REFERENCES venues(id),
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    reason VARCHAR(256) DEFAULT ''
);

CREATE INDEX idx_venue_blocked_venue_id ON venue_blocked_times(venue_id);
CREATE INDEX idx_venue_blocked_time_range ON venue_blocked_times(start_at, end_at);
