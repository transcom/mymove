-- adding column for partial delivery
-- this is stored as jsonb data type, which stores JSON data in binary format
ALTER TABLE mto_shipments
ADD COLUMN partial_deliveries_weight jsonb;

-- Column comments
COMMENT ON COLUMN mto_shipments.partial_deliveries_weight IS 'Partial deliveries defined by weight, stored as JSON';