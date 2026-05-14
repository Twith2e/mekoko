CREATE TABLE cart_items (
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    variant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    unit_price_at_selection BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_cart_item_variant
        FOREIGN KEY(variant_id) REFERENCES product_variants(id) ON DELETE CASCADE,
    CONSTRAINT fk_cart_items_user
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_id_variant_id UNIQUE(user_id, variant_id)
);