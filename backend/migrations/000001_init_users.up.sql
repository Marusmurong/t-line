-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(256),
    nickname VARCHAR(64) DEFAULT '',
    avatar_url VARCHAR(512) DEFAULT '',
    gender SMALLINT DEFAULT 0,
    age SMALLINT DEFAULT 0,
    wechat_openid VARCHAR(128) UNIQUE,
    wechat_unionid VARCHAR(128),
    utr_rating DECIMAL(4,2),
    utr_image VARCHAR(512),
    ball_age SMALLINT DEFAULT 0,
    self_level VARCHAR(20) DEFAULT '',
    member_level SMALLINT DEFAULT 0,
    member_expires_at TIMESTAMPTZ,
    role VARCHAR(20) DEFAULT 'user',
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_wechat_openid ON users(wechat_openid);
CREATE INDEX idx_users_member_level ON users(member_level);
CREATE INDEX idx_users_role ON users(role);

-- Wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    balance DECIMAL(12,2) DEFAULT 0,
    frozen_amount DECIMAL(12,2) DEFAULT 0,
    total_recharged DECIMAL(12,2) DEFAULT 0,
    version INT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Wallet transactions table
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    balance_after DECIMAL(12,2) NOT NULL,
    related_order_id BIGINT,
    remark VARCHAR(256) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_wallet_tx_wallet_id ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_tx_created_at ON wallet_transactions(created_at);

-- Points table
CREATE TABLE IF NOT EXISTS points (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    balance INT DEFAULT 0,
    total_earned INT DEFAULT 0,
    total_spent INT DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
