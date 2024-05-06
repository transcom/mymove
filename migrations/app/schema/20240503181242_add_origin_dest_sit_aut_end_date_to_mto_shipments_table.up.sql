ALTER TABLE mto_shipments
    ADD COLUMN IF NOT EXISTS origin_sit_auth_end_date DATE,
    ADD COLUMN IF NOT EXISTS dest_sit_auth_end_date DATE;

COMMENT ON COLUMN mto_shipments.origin_sit_auth_end_date IS 'receive the initial SIT authorized end date from the origin and subsequent updates';
COMMENT ON COLUMN mto_shipments.dest_sit_auth_end_date IS 'receive the initial SIT authorized end date from the destination and subsequent updates';