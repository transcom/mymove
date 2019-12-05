CREATE TABLE entitlements (
	id uuid PRIMARY KEY NOT NULL,
	dependents_authorized bool,
	total_dependents integer,
	non_temporary_storage bool,
	privately_owned_vehicle bool,
	pro_gear_weight integer,
	pro_gear_weight_spouse integer,
	storage_in_transit integer,
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE move_orders (
	id uuid PRIMARY KEY NOT NULL,
	customer_id uuid REFERENCES customers,
	origin_duty_station_id uuid REFERENCES duty_stations,
	destination_duty_station_id uuid REFERENCES duty_stations,
	entitlements_id uuid REFERENCES entitlements,
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE move_task_orders (
	id uuid PRIMARY KEY NOT NULL,
	move_order_id uuid REFERENCES move_orders,
	reference_id varchar,
	status ghc_approval_status,
	is_available_to_prime bool,
	is_cancelled bool,
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE mto_shipments (
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
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE mto_service_items (
	id uuid PRIMARY KEY NOT NULL,
	move_task_order_id uuid REFERENCES move_task_orders,
	mto_shipment_id uuid REFERENCES mto_shipments,
	meta_id uuid,
	meta_type char,
	created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);