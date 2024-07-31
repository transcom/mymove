CREATE TABLE IF NOT EXISTS mobile_homes (
	id uuid PRIMARY KEY NOT NULL,
    shipment_id    uuid NOT NULL
    CONSTRAINT mobile_home_mto_shipment_id_fkey
    REFERENCES mto_shipments(id),
	make varchar NOT NULL,
	model varchar NOT NULL,
	year int NOT NULL,
	length_in_inches int NOT NULL,
    height_in_inches int NOT NULL,
	widthInInches int NOT NULL,
	created_at timestamp NOT NULL, 
	updated_at timestamp NOT NULL,
	deleted_at timestamp NOT NULL,
);

COMMENT on TABLE mobile_homes IS 'Stores all mobile home shipments, and their details.';
COMMENT on COLUMN mobile_homes.shipment_id IS 'MTO shipment ID associated with this Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.make IS 'Make of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.model IS 'Model of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.year IS 'Year of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.lengthInInches IS 'Length(in) of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.heightInInches IS 'Height(in) of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.widthInInches IS 'Width(in) of the Mobile Home shipment.';
COMMENT on COLUMN mobile_homes.created_at IS 'Date that Mobile Home shipment was created.';
COMMENT on COLUMN mobile_homes.updated_at IS 'Date that Mobile Home shipment was updated.';
COMMENT on COLUMN mobile_homes.deleted_at IS 'Date that the Mobile Home shipment was soft deleted.';

CREATE INDEX mobile_homes_shipment_id_idx ON mobile_homes (shipment_id);
CREATE INDEX mobile_homes_type_idx ON mobile_homes (type);
CREATE INDEX mobile_homes_created_at_idx ON mobile_homes (created_at);
CREATE INDEX mobile_homes_deleted_at_idx ON mobile_homes (deleted_at);