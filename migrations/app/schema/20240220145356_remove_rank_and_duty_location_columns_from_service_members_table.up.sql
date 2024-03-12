-- removing rank & duty location columns from service_members table
-- these values are being pulled from the orders table now
ALTER TABLE service_members
DROP COLUMN rank,
DROP COLUMN duty_location_id;