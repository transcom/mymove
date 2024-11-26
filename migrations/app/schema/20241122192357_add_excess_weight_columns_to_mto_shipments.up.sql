-- add excess_weight_qualified_at and excess_weight_acknowledged_at columns
ALTER TABLE mto_shipments ADD COLUMN IF NOT EXISTS excess_weight_qualified_at timestamp with time zone;
ALTER TABLE mto_shipments ADD COLUMN IF NOT EXISTS excess_weight_acknowledged_at timestamp with time zone;

COMMENT ON COLUMN mto_shipments.excess_weight_qualified_at IS 'The date and time the shipment weight met or exceeded the excess weight qualification threshold';
COMMENT ON COLUMN mto_shipments.excess_weight_acknowledged_at IS 'The date and time the TOO dismissed the risk of excess weight on a shipment.';
