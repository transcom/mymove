-- create mto_service_item_customer_contacts table
-- row is deleted if the mto_service_item is deleted
-- mto service item can only have two customer contacts
CREATE TYPE "customer_contact_type" AS ENUM ('FIRST', 'SECOND');
CREATE TABLE "mto_service_item_customer_contacts" (
    "id" uuid PRIMARY KEY NOT NULL,
    "mto_service_item_id" uuid NOT NULL REFERENCES "mto_service_items" ON DELETE CASCADE,
    "type" "customer_contact_type" NOT NULL,
    "time_military" text NOT NULL,
    "first_available_delivery_date" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT NOW(),
    "updated_at" timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE ("mto_service_item_id", "type")
);

-- comments on columns
COMMENT ON COLUMN "mto_service_item_customer_contacts"."type" IS 'Specify which is the first or second customer contact. Eg. This is the FIRST (type) customer contact with a first available delivery date.';
COMMENT ON COLUMN "mto_service_item_customer_contacts"."time_military" IS 'Specify military time and timezone in text. This is up to the user to decide which format.';
