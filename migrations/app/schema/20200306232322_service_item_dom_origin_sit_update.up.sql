-- Dropping meta concept columns
ALTER TABLE mto_service_items DROP COLUMN meta_id;
ALTER TABLE mto_service_items DROP COLUMN meta_type;

-- Adding additional columns to support Dom. Origin 1st Day SIT
ALTER TABLE mto_service_items
    ADD COLUMN reason             text,
    ADD COLUMN pickup_postal_code text;

-- Column Comment
COMMENT ON COLUMN mto_service_items.reason IS 'Reason for service item.';
COMMENT ON COLUMN mto_service_items.pickup_postal_code IS 'Pickup postal code or zip code for Dom. Origin 1st Day SIT, etc.';