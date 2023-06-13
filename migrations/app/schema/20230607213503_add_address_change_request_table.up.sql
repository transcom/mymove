CREATE TYPE shipment_address_update_status AS enum (
	'REQUESTED',
	'REJECTED',
	'APPROVED'
	);
CREATE TABLE shipment_address_updates
(
	id                                   UUID PRIMARY KEY,
	shipment_id                          uuid REFERENCES mto_shipments(id)                           NOT NULL,
	original_address_id                  uuid                           NOT NULL,
	new_address_id                       uuid                           NOT NULL,
	contractor_remarks                   text                           NOT NULL,
	status                               shipment_address_update_status NOT NULL,
	service_area_changed                 bool                           NOT NULL,
	mileage_bracket_changed              bool                           NOT NULL,
	changed_from_short_haul_to_long_haul bool                           NOT NULL,
	changed_from_long_haul_to_short_haul bool                           NOT NULL,
	office_remarks                       text,
	created_at                           timestamp                      NOT NULL,
	updated_at                           timestamp                      NOT NULL
);

COMMENT ON COLUMN shipment_address_updates.shipment_id IS 'The MTO Shipment ID associated with this address update request';
COMMENT ON COLUMN shipment_address_updates.original_address_id IS 'Original address that was approved for the shipment';
COMMENT ON COLUMN shipment_address_updates.new_address_id IS 'New address being requested';
COMMENT ON COLUMN shipment_address_updates.contractor_remarks IS 'Reason contractor is requesting change to an address that was previously approved';
COMMENT ON COLUMN shipment_address_updates.status IS 'REQUESTED (must be reviewed by TOO), APPROVED (auto-approved, or approved by TOO), or REJECTED (rejected by TOO)';
COMMENT ON COLUMN shipment_address_updates.office_remarks IS 'Remarks from office user who reviewed the request';

CREATE INDEX ON shipment_address_updates (shipment_id);
