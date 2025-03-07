--B22667 Daniel Jordan adding re_country_prn_divisions table
CREATE TABLE IF NOT EXISTS re_country_prn_divisions (
    id          			uuid    	NOT NULL PRIMARY KEY,
    country_id  			uuid    	NOT NULL
    	CONSTRAINT fk_re_country_prn_divisions_re_countries REFERENCES re_countries (id),
    country_prn_dv_id 		text		NOT NULL,
    country_prn_dv_nm  		text    	NOT NULL,
    country_prn_dv_cd		text		NOT NULL,
    command_org_cd			text,
    created_at  			timestamp   NOT NULL DEFAULT NOW(),
    updated_at  			timestamp   NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE re_country_prn_divisions IS 'Stores country principal division data';
COMMENT ON COLUMN re_country_prn_divisions.country_id IS 'The ID for the Country';
COMMENT ON COLUMN re_country_prn_divisions.country_prn_dv_id IS 'The FIPS ID for the country principal division';
COMMENT ON COLUMN re_country_prn_divisions.country_prn_dv_nm IS 'The name of the country principal division';
COMMENT ON COLUMN re_country_prn_divisions.country_prn_dv_cd IS 'The code that represents a country principal division';
COMMENT ON COLUMN re_country_prn_divisions.command_org_cd IS 'The code that represents the responsible command';