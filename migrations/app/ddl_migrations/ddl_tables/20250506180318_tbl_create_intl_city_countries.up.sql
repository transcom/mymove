-- B-23565 Ricky Mettler initial table creation

CREATE TABLE IF NOT EXISTS intl_city_countries (
    id                      uuid        NOT NULL PRIMARY KEY,
    country_id              uuid        NOT NULL
        CONSTRAINT fk_intl_city_countries_re_countries REFERENCES re_countries (id),
    intl_cities_id              uuid        NOT NULL
        CONSTRAINT fk_intl_city_countries_re_intl_cities REFERENCES re_intl_cities (id),
    country_prn_division_id uuid        NOT null
        constraint fk_intl_city_countries_re_country_prn_divisions REFERENCES re_country_prn_divisions (id),
    created_at              timestamp   NOT NULL DEFAULT NOW(),
    updated_at              timestamp   NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE intl_city_countries IS 'Stores association between international country, city and country principal division';
COMMENT ON COLUMN intl_city_countries.country_id IS 'The ID for the Country';
COMMENT ON COLUMN intl_city_countries.intl_cities_id IS 'The ID for the International City';
COMMENT ON COLUMN intl_city_countries.country_prn_division_id IS 'The ID for the country principal division';
