CREATE TABLE move_orders
(
	id uuid PRIMARY KEY NOT NULL,
	customer_id uuid REFERENCES customers,
	origin_duty_station_id uuid REFERENCES duty_stations,
	destination_duty_station_id uuid REFERENCES duty_stations,
	entitlement_id uuid REFERENCES entitlements,
	created_at timestamp WITH TIME ZONE NOT NULL,
	updated_at timestamp WITH TIME ZONE NOT NULL
);

CREATE INDEX ON move_orders (customer_id);
CREATE INDEX ON move_orders (origin_duty_station_id);
CREATE INDEX ON move_orders (destination_duty_station_id);
CREATE INDEX ON move_orders (entitlement_id);