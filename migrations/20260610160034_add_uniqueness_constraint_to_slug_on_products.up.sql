ALTER TABLE products
ALTER COLUMN slug SET NOT NULL;

ALTER TABLE products
ADD CONSTRAINT products_slug_unique UNIQUE (slug);
