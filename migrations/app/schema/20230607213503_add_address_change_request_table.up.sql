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

