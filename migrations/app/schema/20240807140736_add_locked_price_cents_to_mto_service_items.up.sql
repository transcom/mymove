-- Adding locked_price_cents to mto_service_items
ALTER TABLE mto_service_items
ADD COLUMN IF NOT EXISTS locked_price_cents integer;

COMMENT ON COLUMN mto_service_items.locked_price_cents IS 'Locked price for some service items that should never be updated';