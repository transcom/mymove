
CREATE TYPE boat_shipment_type AS enum (
	'HAUL_AWAY',
	'TOW_AWAY'
	);

CREATE TABLE IF NOT EXISTS boat_shipments
(
	id uuid PRIMARY KEY NOT NULL,
	shipment_id    uuid NOT NULL
        CONSTRAINT boat_shipment_mto_shipment_id_fkey
        REFERENCES mto_shipments(id)
        ON DELETE CASCADE,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
    deleted_at timestamp with time zone,
	type boat_shipment_type NOT NULL,
	year integer NOT NULL,
	make varchar NOT NULL,
	model varchar NOT NULL,
	length_in_inches integer NOT NULL,
	width_in_inches integer NOT NULL,
	height_in_inches integer NOT NULL,
	has_Trailer boolean NOT NULL,
	is_roadworthy boolean
);

COMMENT on TABLE boat_shipments IS 'Stores all Boat shipments, and their details.';
COMMENT on COLUMN boat_shipments.shipment_id IS 'MTO shipment ID associated with this Boat shipment.';
COMMENT on COLUMN boat_shipments.type IS 'Type of the Boat shipment.';
COMMENT on COLUMN boat_shipments.year IS 'Year of the Boat.';
COMMENT on COLUMN boat_shipments.make IS 'Make of the Boat.';
COMMENT on COLUMN boat_shipments.model IS 'Mode of the Boat.';
COMMENT on COLUMN boat_shipments.length_in_inches IS 'Length of the Boat in inches.';
COMMENT on COLUMN boat_shipments.width_in_inches IS 'Width of the Boat in inches.';
COMMENT on COLUMN boat_shipments.height_in_inches IS 'Height of the Boat in inches.';
COMMENT on COLUMN boat_shipments.has_Trailer IS 'Does the boat have a trailer.';
COMMENT on COLUMN boat_shipments.is_roadworthy IS 'Is the trailer roadworthy.';

CREATE INDEX boat_shipments_shipment_id_idx ON boat_shipments (shipment_id);
CREATE INDEX boat_shipments_type_idx ON boat_shipments (type);
CREATE INDEX boat_shipments_created_at_idx ON boat_shipments (created_at);
CREATE INDEX boat_shipments_deleted_at_idx ON boat_shipments (deleted_at);