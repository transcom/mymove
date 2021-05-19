ALTER TABLE mto_shipments
    ADD COLUMN counselor_remarks text;
COMMENT ON COLUMN mto_shipments.counselor_remarks IS 'Remarks service counselor has on the MTO Shipment';
