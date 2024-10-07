-- create and set market code
create type market_code_enum as enum ('i', 'd');
ALTER TABLE mto_shipments ADD COLUMN IF NOT EXISTS market_code market_code_enum;
COMMENT ON COLUMN mto_shipments.market_code IS 'Market code indicator for the shipment. i for international and d for destination.';
