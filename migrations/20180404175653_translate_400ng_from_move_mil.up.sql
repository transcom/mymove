-- We need access to a UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Pack rates
SELECT
    uuid_generate_v4() as id,
    schedule,
    LOWER(weight_lbs) as weight_lbs_lower,
    UPPER(weight_lbs) as weight_lbs_upper,
    CAST((rate * 100) as INTEGER) as rate_cents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    INTO tariff400ng_full_pack_rates
    FROM full_packs
;
DROP TABLE full_packs;

-- Unpack rates
SELECT
	uuid_generate_v4() as id,
    schedule,
    CAST((rate * 100000) as INTEGER) as rate_millicents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    INTO tariff400ng_full_unpack_rates
    FROM full_unpacks
;
DROP TABLE full_unpacks;

-- Linehaul
SELECT
	uuid_generate_v4() as id,
    LOWER(dist_mi) as dist_mi_lower,
    UPPER(dist_mi) as dist_mi_upper,
    LOWER(weight_lbs) as weight_lbs_lower,
    UPPER(weight_lbs) as weight_lbs_upper,
    CAST((rate * 100) as INTEGER) as rate_cents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
	CAST(type as TEXT) as type,
    created_at,
    updated_at
    INTO tariff400ng_linehaul_rates
    FROM linehauls
;
DROP TABLE linehauls;

-- Service areas
SELECT
	uuid_generate_v4() as id,
	service_area,
	name,
	services_schedule,
    CAST((linehaul_factor * 100) as INTEGER) as linehaul_factor,
    CAST((orig_dest_service_charge * 100) as INTEGER) as orig_dest_service_charge,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    INTO tariff400ng_service_areas
    FROM service_areas
;
DROP TABLE service_areas;

-- Shorthauls
SELECT
	uuid_generate_v4() as id,
    LOWER(cwt_mi) as cwt_mi_lower,
    UPPER(cwt_mi) as cwt_mi_upper,
    CAST((rate * 100) as INTEGER) as rate_cents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    INTO tariff400ng_shorthaul_rates
    FROM shorthauls
;
DROP TABLE shorthauls;

-- ZIP3s
SELECT
	uuid_generate_v4() as id,
	zip3,
	basepoint_city,
	state,
	service_area,
	rate_area,
	region,
    created_at,
    updated_at
    INTO tariff400ng_zip3s
    FROM zip3s
;
DROP TABLE zip3s;

-- ZIP5 Rate Areas
SELECT
	uuid_generate_v4() as id,
	zip5,
	rate_area,
    created_at,
    updated_at
    INTO tariff400ng_zip5_rate_areas
    FROM zip5_rate_areas
;
DROP TABLE zip5_rate_areas;
