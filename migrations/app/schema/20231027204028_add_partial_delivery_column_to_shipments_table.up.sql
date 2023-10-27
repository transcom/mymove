-- adding column for partial delivery
-- this is stored as an array of integers defining weight to account for potential multiple deliveries out of SIT
ALTER TABLE mto_shipments
ADD COLUMN partial_deliveries_weight integer[];

-- Column comments
COMMENT ON COLUMN mto_shipments.partial_deliveries_weight IS 'Partial deliveries defined by weight, stored as an array of integers';