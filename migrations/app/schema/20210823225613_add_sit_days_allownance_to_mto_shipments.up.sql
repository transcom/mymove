ALTER TABLE mto_shipments ADD COLUMN sit_days_allowance int;

COMMENT ON COLUMN mto_shipments.sit_days_allowance IS 'Total number of SIT days allowed for this shipment, including any sit extensions that have been approved'
