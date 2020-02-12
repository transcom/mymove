CREATE TABLE client_certs (
	id uuid PRIMARY KEY,
	sha256_digest char(64),
	subject text,
	allow_orders_api boolean DEFAULT false,
	allow_air_force_orders_read boolean DEFAULT false,
	allow_air_force_orders_write boolean DEFAULT false,
	allow_army_orders_read boolean DEFAULT false,
	allow_army_orders_write boolean DEFAULT false,
	allow_coast_guard_orders_read boolean DEFAULT false,
	allow_coast_guard_orders_write boolean DEFAULT false,
	allow_marine_corps_orders_read boolean DEFAULT false,
	allow_marine_corps_orders_write boolean DEFAULT false,
	allow_navy_orders_read boolean DEFAULT false,
	allow_navy_orders_write boolean DEFAULT false,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);
