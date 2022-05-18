-- Add a type for the PPM's SIT location.
CREATE TYPE sit_location_type AS ENUM ('ORIGIN','DESTINATION');

COMMENT ON TYPE sit_location_type IS 'The type of location for the PPM''s SIT.';

-- Add new SIT-related columns for PPM shipments.
ALTER TABLE ppm_shipments
	ADD COLUMN sit_location                 sit_location_type,
	ADD COLUMN sit_estimated_weight         integer,
	ADD COLUMN sit_estimated_entry_date     date,
	ADD COLUMN sit_estimated_departure_date date,
	ADD COLUMN sit_estimated_cost           integer;

COMMENT ON COLUMN ppm_shipments.sit_location IS 'Records whether the PPM''s SIT is at the origin or destination.';
COMMENT ON COLUMN ppm_shipments.sit_estimated_weight IS 'The estimated weight of the PPM''s SIT.';
COMMENT ON COLUMN ppm_shipments.sit_estimated_entry_date IS 'The estimated date the PPM''s items will go into SIT.';
COMMENT ON COLUMN ppm_shipments.sit_estimated_departure_date IS 'The estimated date the PPM''s items will come out of SIT.';
COMMENT ON COLUMN ppm_shipments.sit_estimated_cost IS 'The estimated cost (in cents) of the PPM''s SIT.';
