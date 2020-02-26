CREATE TYPE mto_shipment_status AS ENUM (
'SUBMITTED',
'APPROVED',
'REJECTED'
    );

ALTER TABLE mto_shipments ADD COLUMN status mto_shipment_status;

UPDATE mto_shipments SET status = 'SUBMITTED';

ALTER TABLE mto_shipments ALTER COLUMN status SET DEFAULT 'SUBMITTED';
