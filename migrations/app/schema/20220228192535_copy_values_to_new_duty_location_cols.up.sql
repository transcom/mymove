UPDATE service_members
SET duty_location_id = duty_station_id
WHERE duty_location_id IS NULL
  AND duty_station_id IS NOT NULL;

UPDATE orders
SET origin_duty_location_id = origin_duty_station_id
WHERE origin_duty_location_id IS NULL
  AND origin_duty_station_id IS NOT NULL;

UPDATE orders
SET new_duty_location_id = new_duty_station_id
WHERE new_duty_location_id IS NULL;

-- new_duty_location_id is the only one of these fields that is not supposed to be
-- nullable, so now that we've filled it in, we can set that constraint.
-- Note that this blocks writes to the table while it is queued and while it is running.
-- I think this is an acceptable risk, but if we were live we'd probably want to do it
-- more carefully.
ALTER TABLE orders
	ALTER COLUMN new_duty_location_id SET NOT NULL;

COMMENT ON COLUMN orders.origin_duty_location_id IS 'Unique identifier for the duty location the customer is moving from. Not the same as the text version of the name.';
COMMENT ON COLUMN orders.new_duty_location_id IS 'Unique identifier for the duty location the customer is being assigned to. Not the same as the text version of the name.';
COMMENT ON COLUMN service_members.duty_location_id IS 'A foreign key that points to the duty location table - containing the customer''s current duty location';
