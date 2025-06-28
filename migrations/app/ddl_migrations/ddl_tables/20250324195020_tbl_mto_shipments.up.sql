-- B-23373 Brooklyn Welsh - Add actual_gun_safe_weight column to mto_shipments to track gun safe weight across all tickets related to the shipment
ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS termination_comments TEXT,
ADD COLUMN IF NOT EXISTS terminated_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS actual_gun_safe_weight int CHECK (actual_gun_safe_weight >= 0) NULL;

COMMENT ON COLUMN mto_shipments.termination_comments IS 'Comments that give the reason a shipment was terminated for cause';
COMMENT ON COLUMN mto_shipments.terminated_at IS 'The date and time a shipment was terminated for cause';
COMMENT ON COLUMN mto_shipments.actual_gun_safe_weight IS 'Tracks the weight of all gun safe tickets linked to this shipment';