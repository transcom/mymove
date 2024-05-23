ALTER TABLE mto_service_items ADD COLUMN sit_delivery_miles INTEGER DEFAULT NULL;
COMMENT ON COLUMN mto_service_items.sit_delivery_miles IS 'Delivery miles between two SIT addresses';