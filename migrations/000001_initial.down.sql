-- =========================
-- DROP TRIGGERS
-- =========================
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP TRIGGER IF EXISTS products_updated_at ON products;
DROP TRIGGER IF EXISTS orders_updated_at ON orders;
DROP TRIGGER IF EXISTS carts_updated_at ON carts;
DROP TRIGGER IF EXISTS product_images_updated_at ON product_images;

-- =========================
-- DROP TRIGGER FUNCTION
-- =========================
DROP FUNCTION IF EXISTS set_updated_at();

-- =========================
-- DROP INDEXES
-- =========================
DROP INDEX IF EXISTS ux_product_primary_image;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_orders_user_id;

-- =========================
-- DROP TABLES (DEPENDENCY ORDER)
-- =========================
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
-- =========================
-- DROP ENUM TYPES
-- =========================
DROP TYPE IF EXISTS product_status;
DROP TYPE IF EXISTS order_status;
