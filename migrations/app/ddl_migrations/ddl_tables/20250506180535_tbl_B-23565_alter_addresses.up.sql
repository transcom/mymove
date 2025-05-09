-- B-23565 Ricky Mettler adding FK relation to intl_city_countries

ALTER TABLE addresses
ADD COLUMN if not exists intl_city_countries_id	uuid
	CONSTRAINT fk_addresses_intl_city_countries REFERENCES intl_city_countries (id);

COMMENT ON COLUMN addresses.intl_city_countries_id IS 'ID for the associated intl_city_countries record';