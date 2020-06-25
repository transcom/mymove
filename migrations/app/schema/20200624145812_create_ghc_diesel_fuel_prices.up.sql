CREATE TABLE ghc_diesel_fuel_prices
(
	id uuid NOT NULL,
	fuel_price_in_millicents integer NOT NULL,
	publication_date date UNIQUE NOT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);
