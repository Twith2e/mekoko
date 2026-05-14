ALTER TABLE product_variants
ADD created_at TIMESTAMPTZ DEFAULT NOW(),
ADD updated_at TIMESTAMPTZ;