-- Populate the new market_code column for shipments
-- since we do not support OCONUS moves yet, these should all be "d"
UPDATE mto_shipments
SET market_code = 'd'
WHERE market_code IS NULL;

-- fixing typo from previous migration
COMMENT ON COLUMN mto_shipments.market_code IS 'Market code indicator for the shipment. i for international and d for domestic.';