-- B-23575 Brian Manley create country_weekends table to store weekend information for each country

CREATE TABLE IF NOT EXISTS public.country_weekends (
    id          			uuid    	NOT NULL PRIMARY KEY,
    country_id			    uuid		NOT NULL
    	CONSTRAINT fk_country_weekends_re_countries REFERENCES re_countries (id),
	is_monday_weekend		boolean     NOT NULL DEFAULT false,
	is_tuesday_weekend		boolean     NOT NULL DEFAULT false,
	is_wednesday_weekend	boolean     NOT NULL DEFAULT false,
	is_thursday_weekend		boolean     NOT NULL DEFAULT false,
	is_friday_weekend		boolean     NOT NULL DEFAULT false,
	is_saturday_weekend		boolean     NOT NULL DEFAULT false,
	is_sunday_weekend		boolean     NOT NULL DEFAULT false,
    created_at  			timestamp   NOT NULL DEFAULT NOW(),
    updated_at  			timestamp   NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_country_weekends_country_id UNIQUE (country_id)
);

COMMENT ON TABLE country_weekends IS 'Stores designated weekend days associated to a country';
COMMENT ON COLUMN country_weekends.country_id IS 'ID of the country';
COMMENT ON COLUMN country_weekends.is_monday_weekend IS 'Indicates if Monday is a weekend';
COMMENT ON COLUMN country_weekends.is_tuesday_weekend IS 'Indicates if Tuesday is a weekend';
COMMENT ON COLUMN country_weekends.is_wednesday_weekend IS 'Indicates if Wednesday is a weekend';
COMMENT ON COLUMN country_weekends.is_thursday_weekend IS 'Indicates if Thursday is a weekend';
COMMENT ON COLUMN country_weekends.is_friday_weekend IS 'Indicates if Friday is a weekend';
COMMENT ON COLUMN country_weekends.is_saturday_weekend IS 'Indicates if Saturday is a weekend';
COMMENT ON COLUMN country_weekends.is_sunday_weekend IS 'Indicates if Sunday is a weekend';