-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type VARCHAR(40) NOT NULL,
    title VARCHAR(128) NOT NULL,
    content TEXT DEFAULT '',
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    extra_data JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
