CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    base_price BIGINT NOT NULL,
    discount_percentage INT NOT NULL DEFAULT 0 CHECK (discount_percentage >= 0 AND discount_percentage <= 90),
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);