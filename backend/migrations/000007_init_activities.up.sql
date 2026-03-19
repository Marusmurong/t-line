-- Activities table
CREATE TABLE IF NOT EXISTS activities (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    type VARCHAR(20) NOT NULL,
    description TEXT DEFAULT '',
    cover_image VARCHAR(512) DEFAULT '',
    venue_id BIGINT REFERENCES venues(id),
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    registration_deadline TIMESTAMPTZ NOT NULL,
    min_participants INT NOT NULL DEFAULT 0,
    max_participants INT NOT NULL DEFAULT 0,
    current_participants INT NOT NULL DEFAULT 0,
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    level_requirement VARCHAR(20) NOT NULL DEFAULT 'all',
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    cancel_check_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_activities_status ON activities(status);
CREATE INDEX idx_activities_type_status ON activities(type, status);
CREATE INDEX idx_activities_start_at ON activities(start_at);

-- Activity registrations table
CREATE TABLE IF NOT EXISTS activity_registrations (
    id BIGSERIAL PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    order_id BIGINT REFERENCES orders(id),
    status VARCHAR(20) NOT NULL DEFAULT 'registered',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_activity_reg_activity_user ON activity_registrations(activity_id, user_id);
CREATE INDEX idx_activity_reg_user_id ON activity_registrations(user_id);
