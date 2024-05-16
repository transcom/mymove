ALTER TABLE mto_service_items ADD COLUMN IF NOT EXISTS pricing_estimate integer DEFAULT NULL;

COMMENT ON COLUMN mto_service_items.pricing_estimate IS 'This column stores the pricing estimate for a service item.';
