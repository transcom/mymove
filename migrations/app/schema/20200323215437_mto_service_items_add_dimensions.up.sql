CREATE TABLE "mto_service_item_dimensions" (
    "id" uuid PRIMARY KEY NOT NULL,
    "length_thousand_inches" INTEGER NOT NULL,
    "height_thousand_inches" INTEGER NOT NULL,
    "width_thousand_inches" INTEGER NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL);

ALTER TABLE "mto_service_items"
    ADD COLUMN "item_dimension_id" uuid REFERENCES "mto_service_item_dimensions",
    ADD COLUMN "crate_dimension_id" uuid REFERENCES "mto_service_item_dimensions";

-- comments on columns
COMMENT ON COLUMN "mto_service_item_dimensions"."length_thousand_inches" IS 'Length of a crate or item. 1000 thousand inches (thou) = 1 inch.';
COMMENT ON COLUMN "mto_service_item_dimensions"."height_thousand_inches" IS 'Height of a crate or item. 1000 thousand inches (thou) = 1 inch.';
COMMENT ON COLUMN "mto_service_item_dimensions"."width_thousand_inches" IS 'Width of a crate or item. 1000 thousand inches (thou) = 1 inch.';
COMMENT ON COLUMN "mto_service_items"."item_dimension_id" IS 'Dimensions for the item. Ex: Crating a decorated horse head and taking dimensions of the item.';
COMMENT ON COLUMN "mto_service_items"."crate_dimension_id" IS 'Dimensions for the crate. Ex: Crating a decorated horse head and taking dimensions of the create.';