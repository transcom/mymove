ALTER TABLE ppm_shipments
ADD deleted_at timestamp with time zone;

COMMENT ON COLUMN ppm_shipments.deleted_at IS 'Indicates whether the ppm shipment has been soft deleted or not, and when it was soft deleted.';
