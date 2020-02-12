-- Pack rates
INSERT INTO tariff400ng_full_pack_rates
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
    FROM full_packs
;
DROP TABLE full_packs;

-- Unpack rates
INSERT INTO tariff400ng_full_unpack_rates
SELECT
    uuid_generate_v4() as id,
    schedule,
    CAST((rate * 100000) as INTEGER) as rate_millicents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    FROM full_unpacks
;
DROP TABLE full_unpacks;

-- Linehaul
INSERT INTO tariff400ng_linehaul_rates
SELECT
    uuid_generate_v4() as id,
    LOWER(dist_mi) as distance_miles_lower,
    UPPER(dist_mi) as distance_miles_upper,
    LOWER(weight_lbs) as weight_lbs_lower,
    UPPER(weight_lbs) as weight_lbs_upper,
    CAST((rate * 100) as INTEGER) as rate_cents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    CAST(type as TEXT) as type,
    created_at,
    updated_at
    FROM linehauls
;
DROP TABLE linehauls;

-- Service areas
INSERT INTO tariff400ng_service_areas
SELECT
    uuid_generate_v4() as id,
    service_area,
    name,
    services_schedule,
    CAST((linehaul_factor * 100) as INTEGER) as linehaul_factor,
    CAST((orig_dest_service_charge * 100) as INTEGER) as service_charge_cents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    FROM service_areas
;
DROP TABLE service_areas;

-- Shorthauls
INSERT INTO tariff400ng_shorthaul_rates
SELECT
    uuid_generate_v4() as id,
    LOWER(cwt_mi) as cwt_miles_lower,
    UPPER(cwt_mi) as cwt_miles_upper,
    CAST((rate * 100) as INTEGER) as rate_cents,
    LOWER(effective) as effective_date_lower,
    UPPER(effective) as effective_date_upper,
    created_at,
    updated_at
    FROM shorthauls
;
DROP TABLE shorthauls;
