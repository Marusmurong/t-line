-- Products table
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    category VARCHAR(20) NOT NULL,
    sub_category VARCHAR(40) NOT NULL DEFAULT '',
    name VARCHAR(128) NOT NULL,
    description TEXT DEFAULT '',
    cover_image VARCHAR(512) DEFAULT '',
    images JSONB DEFAULT '[]',
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    original_price DECIMAL(10,2) NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    sales_count INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 0,
    attributes JSONB DEFAULT '{}',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_products_category_status ON products(category, status);
CREATE INDEX idx_products_sort ON products(sort_order, id);

-- Product SKUs table
CREATE TABLE IF NOT EXISTS product_skus (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL DEFAULT '',
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    attributes JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_product_skus_product_id ON product_skus(product_id);
