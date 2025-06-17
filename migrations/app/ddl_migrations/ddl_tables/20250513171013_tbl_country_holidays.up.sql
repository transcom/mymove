-- B-23575 Brian Manley create country_holidays table to store holiday information for each country

CREATE TABLE IF NOT EXISTS public.country_holidays (
    id                  uuid    	NOT NULL PRIMARY KEY,
    country_id			uuid		NOT NULL
    	CONSTRAINT fk_country_holidays_re_countries REFERENCES re_countries (id),
    holiday_name 		text		NOT NULL,
    observation_date	date        NOT NULL,
    created_at  		timestamp   NOT NULL DEFAULT NOW(),
    updated_at  		timestamp   NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_country_holidays_key UNIQUE (country_id, holiday_name, observation_date)
);

CREATE INDEX IF NOT EXISTS country_holidays_country_id_idx ON country_holidays(country_id);

COMMENT ON TABLE country_holidays IS 'Stores holidays associated to a country with the observation date';
COMMENT ON COLUMN country_holidays.country_id IS 'ID of the country';
COMMENT ON COLUMN country_holidays.holiday_name IS 'Name of the holiday';
COMMENT ON COLUMN country_holidays.observation_date IS 'Observation date for the holiday';