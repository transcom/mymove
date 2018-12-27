CREATE TABLE client_certs (
    sha256_digest char(64) PRIMARY KEY,
    subject text,
	allow_dps_auth_api boolean DEFAULT false,
	allow_orders_api boolean DEFAULT false,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);
