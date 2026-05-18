CREATE TABLE orders(
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    subtotal BIGINT NOT NULL,
    total_amount BIGINT NOT NULL,
    delivery_fee BIGINT NOT NULL,
    currency TEXT NOT NULL,
    delivery_status TEXT NOT NULL,
    payment_status TEXT NOT NULL,
    discount_amount BIGINT,
    ordered_at TIMESTAMPTZ NOT NULL,
    delivered_at TIMESTAMPTZ,
    shipping_address_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_shipping_address FOREIGN KEY(shipping_address_id) REFERENCES shipping_addresses(id)
);