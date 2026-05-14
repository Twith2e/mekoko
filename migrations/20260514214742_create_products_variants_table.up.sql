CREATE TABLE product_variants (
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    product_id BIGINT NOT NULL,
    size TEXT,
    color TEXT NOT NULL,
    stock_quantity BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_product_variants_product
        FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT uq_product_variants_product_color UNIQUE (product_id, color)
);