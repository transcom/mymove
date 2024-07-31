CREATE TABLE IF NOT EXISTS mobile_home (
	id uuid PRIMARY KEY NOT NULL,
    shipment_id    uuid NOT NULL
    CONSTRAINT mobile_home_mto_shipment_id_fkey
    REFERENCES mto_shipments(id),
	make varchar NOT NULL,
	model varchar NOT NULL,
	mh_year int NOT NULL,
	mh_length int NOT NULL,
    height int NOT NULL,
	width int NOT NULL,
	created_at timestamp NOT NULL, 
	updated_at timestamp NOT NULL,
	deleted_at timestamp NOT NULL,
);

COMMENT on TABLE mobile_home IS 'Stores all mobile home shipments, and their details.';
COMMENT on COLUMN mobile_home.shipment_id IS 'MTO shipment ID associated with this Mobile Home shipment.';
COMMENT on COLUMN mobile_home.make IS 'Make of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.model IS 'Model of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.mh_year IS 'Year of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.mh_length IS 'Length of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.height IS 'Height of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.width IS 'Width of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.created_at IS 'Date that Mobile Home shipment was created.';
COMMENT on COLUMN mobile_home.updated_at IS 'Date that Mobile Home shipment was updated.';
COMMENT on COLUMN mobile_home.deleted_at IS 'Date that the Mobile Home shipment was soft deleted.';