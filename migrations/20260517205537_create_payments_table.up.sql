CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    order_id BIGINT NOT NULL,
    provider TEXT NOT NULL,
    provider_reference TEXT NOT NULL,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders(id)
);