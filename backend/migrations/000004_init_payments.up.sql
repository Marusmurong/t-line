-- Payments table
CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    payment_no VARCHAR(64) NOT NULL UNIQUE,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    method VARCHAR(20) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    balance_amount DECIMAL(12,2) DEFAULT 0,
    wechat_amount DECIMAL(12,2) DEFAULT 0,
    wechat_trade_no VARCHAR(64),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_payment_no ON payments(payment_no);
CREATE INDEX idx_payments_status ON payments(status);

-- Coupons table
CREATE TABLE IF NOT EXISTS coupons (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(20) NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    min_amount DECIMAL(10,2) DEFAULT 0,
    applicable_types JSONB DEFAULT '[]',
    total_count INT NOT NULL,
    used_count INT DEFAULT 0,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    status SMALLINT DEFAULT 1
);

CREATE INDEX idx_coupons_status ON coupons(status);

-- User coupons table
CREATE TABLE IF NOT EXISTS user_coupons (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    coupon_id BIGINT NOT NULL REFERENCES coupons(id),
    status SMALLINT DEFAULT 0,
    used_order_id BIGINT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_coupons_user_id ON user_coupons(user_id);
CREATE INDEX idx_user_coupons_coupon_id ON user_coupons(coupon_id);
CREATE INDEX idx_user_coupons_status ON user_coupons(status);
