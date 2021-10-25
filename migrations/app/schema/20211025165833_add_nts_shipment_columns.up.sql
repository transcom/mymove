CREATE TABLE storage_facilities (
	id uuid
		CONSTRAINT storage_facilities_pkey PRIMARY KEY,
	created_at timestamp not null,
    updated_at timestamp not null,
	facility_name varchar,
	address_id uuid
		CONSTRAINT storage_facilities_address_id_fkey REFERENCES addresses,
	lot_number varchar(255),
	phone varchar(255),
	email varchar(255)
);

CREATE INDEX ON storage_facilities (address_id);

ALTER TABLE orders
	ADD COLUMN nts_tac varchar(255),
	ADD COLUMN nts_sac varchar(255);

ALTER TABLE mto_shipments
	ADD COLUMN external_vendor boolean DEFAULT false,
	ADD COLUMN storage_facility_id uuid
		CONSTRAINT mto_shipments_storage_facility_id REFERENCES storage_facilities,
	ADD COLUMN service_order_number varchar;

CREATE INDEX ON mto_shipments (storage_facility_id);

COMMENT ON TABLE storage_facilities IS 'Storage facilities for NTS and NTS-Release shipments';
COMMENT ON COLUMN storage_facilities.created_at IS 'Date & time the storage facility was created';
COMMENT ON COLUMN storage_facilities.updated_at IS 'Date & time the storage facility was updated';
COMMENT ON COLUMN storage_facilities.facility_name IS 'Name of storage facility';
COMMENT ON COLUMN storage_facilities.address_id IS 'The address of the storage facility';
COMMENT ON COLUMN storage_facilities.lot_number IS 'Lot number where goods are stored within the storage facility';
COMMENT ON COLUMN storage_facilities.phone IS 'Phone number for contacting storage facility';
COMMENT ON COLUMN storage_facilities.email IS 'Email address for contacting storage facility';

COMMENT ON COLUMN orders.tac IS '(For HHG shipments) Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Transportation Ordering Officer adds this information to the MTO.'
COMMENT ON COLUMN orders.sac IS '(For HHG shipments) Shipment Account Classification - used for accounting';
COMMENT ON COLUMN orders.nts_tac IS '(For NTS shipments) Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Transportation Ordering Officer adds this information to the MTO.'
COMMENT ON COLUMN orders.nts_sac IS '(For NTS shipments) Shipment Account Classification - used for accounting';

COMMENT ON COLUMN mto_shipments.external_vendor IS 'Whether this shipment is handled by an external vendor, or by the prime'
COMMENT ON COLUMN mto_shipments.storage_facility_id IS 'The storage facility for an NTS shipment where items are stored'
COMMENT ON COLUMN mto_shipments.service_order_number IS 'The order number for a shipment in TOPS'
