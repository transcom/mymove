CREATE TYPE mto_shipment_type AS ENUM (
'HHG',
'INTERNATIONAL_HHG',
'INTERNATIONAL_UB'
    );

ALTER TABLE mto_shipments
	ADD COLUMN shipment_type mto_shipment_type NOT NULL DEFAULT 'HHG';
