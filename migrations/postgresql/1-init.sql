-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE "nexmedis_users" (
    "id" TEXT NOT NULL,
    "role_id" INTEGER NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "last_login" TIMESTAMPTZ(3),
    "balance" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_users_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "nexmedis_users_email_key" ON "nexmedis_users"("email");
CREATE INDEX "nexmedis_users_deleted_at_last_login_idx" ON "nexmedis_users"("deleted_at", "last_login");

CREATE TABLE "nexmedis_user_addresses" (
    "id" SERIAL NOT NULL,
    "user_id" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "city" TEXT NOT NULL,
    "province" TEXT NOT NULL,
    "latitude" TEXT,
    "longitude" TEXT,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_user_addresses_pkey" PRIMARY KEY ("id")
);

CREATE INDEX "nexmedis_user_addresses_user_id_deleted_at_idx" ON "nexmedis_user_addresses"("user_id", "deleted_at");

CREATE TABLE "nexmedis_user_carts" (
    "id" SERIAL NOT NULL,
    "user_id" TEXT NOT NULL,
    "product_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL DEFAULT 1,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_user_carts_pkey" PRIMARY KEY ("id")
);

CREATE INDEX "nexmedis_user_carts_user_id_product_id_idx" ON "nexmedis_user_carts"("user_id", "product_id");

CREATE TABLE "nexmedis_user_transactions" (
    "id" SERIAL NOT NULL,
    "user_id" TEXT NOT NULL,
    "transaction_status_id" INTEGER NOT NULL,
    "total_amount" DOUBLE PRECISION NOT NULL,
    "invoice_number" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_user_transactions_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "nexmedis_user_transactions_invoice_number_key" ON "nexmedis_user_transactions"("invoice_number");
CREATE INDEX "nexmedis_user_transactions_user_id_transaction_status_id_idx" ON "nexmedis_user_transactions"("user_id", "transaction_status_id", "updated_at");

CREATE TABLE "nexmedis_products" (
    "id" SERIAL NOT NULL,
    "sku" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "slug" TEXT NOT NULL,
    "description" TEXT,
    "color" TEXT,
    "size" TEXT,
    "price" DOUBLE PRECISION NOT NULL,
    "stock" INTEGER NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_products_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "nexmedis_products_sku_key" ON "nexmedis_products"("sku");
CREATE UNIQUE INDEX "nexmedis_products_name_key" ON "nexmedis_products"("name");
CREATE UNIQUE INDEX "nexmedis_products_slug_key" ON "nexmedis_products"("slug");

CREATE TABLE "nexmedis_master_user_roles" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_master_user_roles_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "nexmedis_master_user_roles_name_key" ON "nexmedis_master_user_roles"("name");

CREATE TABLE "nexmedis_master_transaction_statuses" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ(3) NOT NULL,
    "updated_by" TEXT NOT NULL,
    "deleted_at" TIMESTAMPTZ(3),
    "deleted_by" TEXT,

    CONSTRAINT "nexmedis_master_transaction_statuses_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "nexmedis_master_transaction_statuses_name_key" ON "nexmedis_master_transaction_statuses"("name");

ALTER TABLE "nexmedis_users" ADD CONSTRAINT "nexmedis_users_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "nexmedis_master_user_roles"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "nexmedis_user_addresses" ADD CONSTRAINT "nexmedis_user_addresses_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "nexmedis_users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "nexmedis_user_carts" ADD CONSTRAINT "nexmedis_user_carts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "nexmedis_users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "nexmedis_user_carts" ADD CONSTRAINT "nexmedis_user_carts_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "nexmedis_products"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "nexmedis_user_transactions" ADD CONSTRAINT "nexmedis_user_transactions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "nexmedis_users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "nexmedis_user_transactions" ADD CONSTRAINT "nexmedis_user_transactions_transaction_status_id_fkey" FOREIGN KEY ("transaction_status_id") REFERENCES "nexmedis_master_transaction_statuses"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

INSERT INTO "nexmedis_master_user_roles" ("name", "created_by", "updated_by", "updated_at")
VALUES 
  ('admin', 'seed', 'seed', CURRENT_TIMESTAMP),
  ('user',  'seed', 'seed', CURRENT_TIMESTAMP);

INSERT INTO "nexmedis_master_transaction_statuses" ("name", "created_by", "updated_by", "updated_at")
VALUES 
  ('complete', 'seed', 'seed', CURRENT_TIMESTAMP),
  ('failed',   'seed', 'seed', CURRENT_TIMESTAMP);

INSERT INTO "nexmedis_users" (
  "id", 
  "role_id", 
  "email", 
  "password", 
  "created_by", 
  "updated_by", 
  "updated_at"
)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  2,
  'user@example.com',
  crypt('userPassword', gen_salt('bf')),
  'seed',
  'seed',
  CURRENT_TIMESTAMP
);

INSERT INTO "nexmedis_user_addresses" (
  "user_id", 
  "address", 
  "city", 
  "province", 
  "created_by", 
  "updated_by", 
  "updated_at"
)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  'Jl. Kalijudan 2',
  'Surabaya',
  'Jawa Timur',
  'seed',
  'seed',
  CURRENT_TIMESTAMP
);

INSERT INTO "nexmedis_products" (
  "sku", 
  "name", 
  "slug", 
  "description", 
  "color", 
  "size", 
  "price", 
  "stock", 
  "created_by", 
  "updated_by", 
  "updated_at"
)
VALUES 
  ('SKU001', 'Product 1', 'product-1', 'Description for product 1', 'red',  'm', 20000, 100, 'seed', 'seed', CURRENT_TIMESTAMP),
  ('SKU002', 'Product 2', 'product-2', 'Description for product 2', 'green',  's', 15000, 30, 'seed', 'seed', CURRENT_TIMESTAMP),
  ('SKU003', 'Product 3', 'product-3', 'Description for product 3', 'yellow',  'xl', 30000, 20, 'seed', 'seed', CURRENT_TIMESTAMP),
  ('SKU004', 'Product 4', 'product-4', 'Description for product 4', 'cyan',  'xs', 10000, 10, 'seed', 'seed', CURRENT_TIMESTAMP),
  ('SKU005', 'Product 5', 'product-5', 'Description for product 5', 'orange',  'xxl', 35000, 10, 'seed', 'seed', CURRENT_TIMESTAMP),
  ('SKU006', 'Product 6', 'product-6', 'Description for product 6', 'blue', 'l', 25000, 50,  'seed', 'seed', CURRENT_TIMESTAMP);

INSERT INTO "nexmedis_user_carts" (
  "user_id", 
  "product_id", 
  "quantity", 
  "created_by", 
  "updated_by", 
  "updated_at"
)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  1,
  2,
  'seed',
  'seed',
  CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS nexmedis_user_carts;
DROP TABLE IF EXISTS nexmedis_user_addresses;
DROP TABLE IF EXISTS nexmedis_user_transactions;
DROP TABLE IF EXISTS nexmedis_products;
DROP TABLE IF EXISTS nexmedis_users;
DROP TABLE IF EXISTS nexmedis_master_user_roles;
DROP TABLE IF EXISTS nexmedis_master_transaction_statuses;
