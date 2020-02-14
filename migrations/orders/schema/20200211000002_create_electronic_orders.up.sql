CREATE TABLE electronic_orders (
	id uuid PRIMARY KEY,
	orders_number text,
	edipi text,
	issuer text,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

CREATE UNIQUE INDEX electronic_orders_index_by_issuer_and_orders_number ON electronic_orders (issuer, orders_number);
CREATE INDEX electronic_orders_index_by_issuer_and_edipi ON electronic_orders (issuer, edipi);
