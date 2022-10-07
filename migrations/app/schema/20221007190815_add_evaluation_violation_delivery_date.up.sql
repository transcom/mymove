ALTER TABLE evaluation_reports 
  ADD COLUMN observed_delivery_date date;

COMMENT ON COLUMN evaluation_reports.observed_delivery_date IS 'Observed delivery date';