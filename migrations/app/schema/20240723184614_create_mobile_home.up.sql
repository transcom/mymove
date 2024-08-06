CREATE TABLE IF NOT EXISTS mobile_homes (
	id uuid PRIMARY KEY NOT NULL,
    shipment_id  uuid NOT NULL
    CONSTRAINT mobile_home_mto_shipment_id_fkey
    REFERENCES mto_shipments(id)
	ON DELETE CASCADE,
	make varchar NOT NULL,
	model varchar NOT NULL,
	year int NOT NULL,
	length_in_inches int NOT NULL,
    height_in_inches int NOT NULL,
	width_in_inches int NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamp with time zone
);

COMMENT on TABLE mobile_homes IS 'Stores all mobile home shipments, and their details.';
COMMENT on COLUMN mobile_homes.shipment_id IS 'MTO shipment ID associated with this Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.make IS 'Make of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.model IS 'Model of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.year IS 'Year of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.length_in_inches IS 'Length(in) of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.height_in_inches IS 'Height(in) of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.width_in_inches IS 'Width(in) of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.created_at IS 'Date that Mobile Home shipment was created.';
COMMENT on COLUMN mobile_homes.updated_at IS 'Date that Mobile Home shipment was updated.';
COMMENT on COLUMN mobile_homes.deleted_at IS 'Date that the Mobile Home shipment was soft deleted.';

CREATE INDEX mobile_homes_shipment_id_idx ON mobile_homes (shipment_id);
CREATE INDEX mobile_homes_deleted_at_idx ON mobile_homes (deleted_at);