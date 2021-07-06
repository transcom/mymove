ALTER TABLE mto_shipments ADD COLUMN deleted_at timestamp with time zone;
-- Column comment
COMMENT ON COLUMN mto_shipments.deleted_at IS 'Indicates whether the shipment has been soft deleted or not, and when it was soft deleted.';
