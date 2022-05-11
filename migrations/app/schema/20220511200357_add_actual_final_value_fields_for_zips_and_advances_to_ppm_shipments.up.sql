-- Add new columns to ppm_shipments table to track "actual" values once a customer has moved

ALTER TABLE ppm_shipments
	ADD COLUMN actual_pickup_postal_code varchar,
	ADD COLUMN actual_destination_postal_code varchar,
	ADD COLUMN advance_received bool,
	ADD COLUMN actual_advance int;

comment on column ppm_shipments.actual_pickup_postal_code is 'Tracks the actual postal code where the PPM shipment began.';
comment on column ppm_shipments.actual_destination_postal_code is 'Tracks the actual destination postal code for PPM shipment once customer moved the shipment.';
comment on column ppm_shipments.advance_received is 'Indicates if a customer actually received an advance for their PPM shipment.';
comment on column ppm_shipments.actual_advance is 'Tracks the amount a customer received for their advance; amount should be a percentage of estimated incentive.';
