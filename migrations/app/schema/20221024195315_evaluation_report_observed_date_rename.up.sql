ALTER TABLE evaluation_reports
  DROP COLUMN observed_date,
  ADD COLUMN observed_shipment_delivery_date date,
  ADD COLUMN observed_shipment_physical_pickup_date date;


COMMENT ON COLUMN evaluation_reports.observed_shipment_delivery_date IS 'Indicates shipment delivery date was different from scheduled';
COMMENT ON COLUMN evaluation_reports.observed_shipment_physical_pickup_date IS 'Indicates shipment pickup date was different from scheduled';
