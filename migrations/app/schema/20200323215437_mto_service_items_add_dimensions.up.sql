-- create dimensions table
-- row is deleted if the mto_service_item is deleted
-- mto service item can only have an item or crate dimension
CREATE TYPE "dimension_type" AS ENUM ('ITEM', 'CRATE');
CREATE TABLE "mto_service_item_dimensions" (
    "id" uuid PRIMARY KEY NOT NULL,
    "mto_service_item_id" uuid NOT NULL REFERENCES "mto_service_items" ON DELETE CASCADE,
    "type" "dimension_type" NOT NULL,
    "length_thousandth_inches" INTEGER NOT NULL,
    "height_thousandth_inches" INTEGER NOT NULL,
    "width_thousandth_inches" INTEGER NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    UNIQUE ("mto_service_item_id", "type")
);

-- comments on columns
COMMENT ON COLUMN "mto_service_item_dimensions"."type" IS 'Specify what is being measured, the item itself or crate. Eg. Measuring dimensions for a decorated horse head, ITEM type.';
COMMENT ON COLUMN "mto_service_item_dimensions"."length_thousandth_inches" IS 'Length of a crate or item. 1000 thousandth inches (thou) = 1 inch.';
COMMENT ON COLUMN "mto_service_item_dimensions"."height_thousandth_inches" IS 'Height of a crate or item. 1000 thousandth inches (thou) = 1 inch.';
COMMENT ON COLUMN "mto_service_item_dimensions"."width_thousandth_inches" IS 'Width of a crate or item. 1000 thousandth inches (thou) = 1 inch.';