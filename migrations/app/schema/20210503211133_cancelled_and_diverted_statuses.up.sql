ALTER TYPE mto_shipment_status ADD VALUE 'CANCELED';
ALTER TYPE mto_shipment_status ADD VALUE 'DIVERSION_REQUESTED';
COMMENT ON COLUMN mto_shipments.status IS 'The status of a shipment.';
ALTER TABLE mto_shipments
    ADD COLUMN diversion BOOLEAN NOT NULL DEFAULT FALSE;
COMMENT ON COLUMN mto_shipments.diversion IS 'Indicate if the shipment is part of a diversion';
