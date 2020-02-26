CREATE TABLE mto_shipments
(
	id uuid PRIMARY KEY NOT NULL,
	move_task_order_id uuid REFERENCES move_task_orders,
	scheduled_pickup_date date,
	requested_pickup_date date,
	customer_remarks text,
	pickup_address_id uuid REFERENCES addresses,
	destination_address_id uuid REFERENCES addresses,
	secondary_pickup_address_id uuid REFERENCES addresses,
	secondary_delivery_address_id uuid REFERENCES addresses,
	prime_estimated_weight integer,
	prime_estimated_weight_recorded_date date,
	prime_actual_weight integer,
	created_at timestamp WITH TIME ZONE NOT NULL,
	updated_at timestamp WITH TIME ZONE NOT NULL
);

CREATE INDEX ON mto_shipments (move_task_order_id);
CREATE INDEX ON mto_shipments (pickup_address_id);
CREATE INDEX ON mto_shipments (destination_address_id);
CREATE INDEX ON mto_shipments (secondary_pickup_address_id);
CREATE INDEX ON mto_shipments (secondary_delivery_address_id);