ALTER TABLE mto_shipments ADD COLUMN calculated_billable_weight int;

COMMENT ON COLUMN mto_shipments.calculated_billable_weight IS 'The weight that has been calculated by the system for a MTO shipment';
