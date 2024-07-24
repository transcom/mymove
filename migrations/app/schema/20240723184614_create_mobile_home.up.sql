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
	requested_pickup_date date,
    requested_delivery_date date,
	pickup_address varchar NOT NULL,
	destination_address varchar,
    origin_address varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamptz,
	has_secondary_pickup_address boolean,
	has_secondary_destination_address boolean,
	secondary_pickup_address varchar,
    secondary_destination_address varchar,
	receiving_agent varchar,
	counselor_remarks varchar,
	customer_remarks varchar
);

COMMENT on TABLE mobile_home IS 'Stores all mobile home shipments, and their details.';
COMMENT on COLUMN mobile_home.shipment_id IS 'MTO shipment ID associated with this PPM shipment.';
COMMENT on COLUMN mobile_home.make IS 'Make of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.model IS 'Model of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.mh_year IS 'Year of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.mh_length IS 'Length of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.height IS 'Height of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.width IS 'Width of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.requested_pickup_date IS 'Requested date of the Mobile Home shipment pickup by prime.';
COMMENT on COLUMN mobile_home.requested_delivery_date IS 'Requested date of the Mobile Home shipment delivery by prime.';
COMMENT on COLUMN mobile_home.pickup_address IS 'Address of where the prime can pickup the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.destination_address IS 'Address of where the prime can deliver the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.origin_address IS 'Origin address of the of the service member.';
COMMENT on COLUMN mobile_home.secondary_pickup_address IS 'Secondary/Backup pickup address of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.has_secondary_pickup_address IS 'true/false if Mobile Home Shipment has a secondary pickup address.';
COMMENT on COLUMN mobile_home.has_secondary_destination_address IS 'true/false if Mobile Home Shipment has a secondary destination address.';
COMMENT on COLUMN mobile_home.secondary_destination_address IS 'Secondary/Backup destination address of the Mobile Home shipment.';
COMMENT on COLUMN mobile_home.receiving_agent IS 'Receiving agent of the Mobile Home Shipment.';
COMMENT on COLUMN mobile_home.counselor_remarks IS 'Counselor comments on Mobile Home Shipment.';
COMMENT on COLUMN mobile_home.customer_remarks IS 'Customer comments on Mobile Home Shipment.';
