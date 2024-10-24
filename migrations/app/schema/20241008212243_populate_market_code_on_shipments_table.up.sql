-- Set temp timeout due to potentially large modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

-- Populate the new market_code column for shipments
-- since we do not support OCONUS moves yet, these should all be "d"
UPDATE mto_shipments
SET market_code = 'd'
WHERE market_code IS NULL;

-- Add a NOT NULL constraint to the market_code column after we populate with data
ALTER TABLE mto_shipments
ALTER COLUMN market_code SET NOT NULL;

-- fixing typo from previous migration
COMMENT ON COLUMN mto_shipments.market_code IS 'Market code indicator for the shipment. i for international and d for domestic.';