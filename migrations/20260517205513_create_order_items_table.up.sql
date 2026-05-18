CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    variant_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    unit_price BIGINT NOT NULL,
    total_price BIGINT NOT NULL,
    product_name_snapshot TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id),
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_variant FOREIGN KEY(variant_id) REFERENCES product_variants(id)
)