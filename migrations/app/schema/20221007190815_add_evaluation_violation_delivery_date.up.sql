ALTER TABLE evaluation_reports 
  ADD COLUMN observed_delivery_date date;

COMMENT ON COLUMN evaluation_reports.observed_delivery_date IS 'Observed delivery date';

UPDATE pws_violations
SET additional_data_elem = 'observedClaimsResponseDate'
WHERE id = '1261c17d-5229-4004-a17c-ed7765c7d491';