ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS termination_comments TEXT,
ADD COLUMN IF NOT EXISTS terminated_at TIMESTAMP;

COMMENT ON COLUMN mto_shipments.termination_comments IS 'Comments that give the reason a shipment was terminated for cause';
COMMENT ON COLUMN mto_shipments.terminated_at IS 'The date and time a shipment was terminated for cause';