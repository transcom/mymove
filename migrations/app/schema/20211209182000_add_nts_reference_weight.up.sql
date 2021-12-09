alter table mto_shipments
    add column nts_recorded_weight integer;

COMMENT ON COLUMN mto_shipments.nts_recorded_weight IS 'Previously recorded weight used for NTS shipment';
