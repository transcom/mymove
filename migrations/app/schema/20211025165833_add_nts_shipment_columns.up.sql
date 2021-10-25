CREATE TABLE storage_facilities (
	id uuid PRIMARY KEY NOT NULL,
	name varchar,
	address_id uuid REFERENCES addresses,
	lot_number varchar,
	phone varchar,
	email varchar
);

CREATE INDEX ON storage_facilities (address_id);

ALTER TABLE orders
	ADD COLUMN nts_tac varchar,
	ADD COLUMN nts_sac varchar;

ALTER TABLE mto_shipments
	ADD COLUMN external_vendor boolean,
	ADD COLUMN storage_facility_id uuid REFERENCES storage_facilities,
	ADD COLUMN service_order_number varchar;

CREATE INDEX ON mto_shipments (storage_facility_id);
