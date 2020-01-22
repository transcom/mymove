CREATE TYPE mto_shipment_status AS ENUM (
'SUBMITTED',
'APPROVED',
'REJECTED'
    );

ALTER TABLE mto_shipments
	ADD COLUMN status mto_shipment_status NOT NULL DEFAULT 'SUBMITTED';
