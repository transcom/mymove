-- B-23565 Ricky Mettler initial table creation

CREATE TABLE IF NOT EXISTS re_intl_cities (
    id          			uuid    	NOT NULL PRIMARY KEY,
    city_name				text		NOT NULL,
    country_id  			uuid    	NOT NULL
    	CONSTRAINT fk_re_intl_cities_re_countries REFERENCES re_countries (id),
    created_at  			timestamp   NOT NULL DEFAULT NOW(),
    updated_at  			timestamp   NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE re_intl_cities IS 'Stores international cities';
COMMENT ON COLUMN re_intl_cities.city_name IS 'The name of the City';
COMMENT ON COLUMN re_intl_cities.country_id IS 'The ID for the Country';