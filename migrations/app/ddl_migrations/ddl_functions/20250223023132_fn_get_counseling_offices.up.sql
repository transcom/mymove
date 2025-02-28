--B-22660 Daniel Jordan added get_duty_location_info
CREATE OR REPLACE FUNCTION get_duty_location_info(p_duty_location_id UUID)
RETURNS TABLE (duty_addr_id UUID, is_oconus BOOLEAN)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT dl.address_id, a.is_oconus
    FROM duty_locations dl
    JOIN addresses a ON a.id = dl.address_id
    WHERE dl.id = p_duty_location_id;
END;
$$;

--B-22660 Daniel Jordan added get_service_affiliation
CREATE OR REPLACE FUNCTION get_service_affiliation(p_service_member_id UUID)
RETURNS TEXT
LANGUAGE plpgsql AS $$
DECLARE
    service_affiliation TEXT;
BEGIN
    SELECT affiliation INTO service_affiliation
    FROM service_members
    WHERE id = p_service_member_id;

    RETURN service_affiliation;
END;
$$;

--B-22660 Daniel Jordan added get_department_indicator
CREATE OR REPLACE FUNCTION get_department_indicator(p_service_affiliation TEXT)
RETURNS TEXT
LANGUAGE plpgsql AS $$
DECLARE
    dept_indicator TEXT;
BEGIN
    IF p_service_affiliation IN ('AIR_FORCE', 'SPACE_FORCE') THEN
        dept_indicator := 'AIR_AND_SPACE_FORCE';
    ELSIF p_service_affiliation IN ('NAVY', 'MARINES') THEN
        dept_indicator := 'NAVY_AND_MARINES';
    ELSIF p_service_affiliation = 'ARMY' THEN
        dept_indicator := 'ARMY';
    ELSIF p_service_affiliation = 'COAST_GUARD' THEN
        dept_indicator := 'COAST_GUARD';
    ELSE
        RAISE EXCEPTION 'Invalid affiliation: %', p_service_affiliation;
    END IF;

    RETURN dept_indicator;
END;
$$;

--B-22660 Daniel Jordan added get_gbloc_indicator
CREATE OR REPLACE FUNCTION get_gbloc_indicator(p_duty_addr_id UUID, p_dept_indicator TEXT)
RETURNS TEXT
LANGUAGE plpgsql AS $$
DECLARE
    gbloc_indicator TEXT;
BEGIN
    SELECT j.code INTO gbloc_indicator
    FROM addresses a
    JOIN v_locations v ON a.us_post_region_cities_id = v.uprc_id
    JOIN re_oconus_rate_areas o ON v.uprc_id = o.us_post_region_cities_id
    JOIN re_rate_areas r ON o.rate_area_id = r.id
    JOIN gbloc_aors g ON o.id = g.oconus_rate_area_id
    JOIN jppso_regions j ON g.jppso_regions_id = j.id
    WHERE a.id = p_duty_addr_id
        AND (g.department_indicator = p_dept_indicator OR g.department_indicator IS NULL)
    LIMIT 1;

    IF gbloc_indicator IS NULL THEN
        RAISE EXCEPTION 'Cannot determine GBLOC for duty location';
    END IF;

    RETURN gbloc_indicator;
END;
$$;

--B-22660 Daniel Jordan added fetch_counseling_offices_for_oconus
CREATE OR REPLACE FUNCTION fetch_counseling_offices_for_oconus(p_duty_location_id UUID, p_gbloc_indicator TEXT)
RETURNS TABLE (id UUID, name TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT toff.id, toff.name
    FROM duty_locations dl
    JOIN addresses a ON dl.address_id = a.id
    JOIN v_locations v ON (a.us_post_region_cities_id = v.uprc_id OR v.uprc_id IS NULL)
         AND a.country_id = v.country_id
    JOIN re_oconus_rate_areas r ON r.us_post_region_cities_id = v.uprc_id
    JOIN gbloc_aors ga ON ga.oconus_rate_area_id = r.id
    JOIN jppso_regions j ON ga.jppso_regions_id = j.id
    JOIN transportation_offices toff ON j.code = toff.gbloc
    JOIN addresses toff_addr ON toff.address_id = toff_addr.id
    LEFT JOIN zip3_distances zd
      ON (
         (substring(a.postal_code, 1, 3) = zd.from_zip3 AND substring(toff_addr.postal_code, 1, 3) = zd.to_zip3)
         OR
         (substring(a.postal_code, 1, 3) = zd.to_zip3 AND substring(toff_addr.postal_code, 1, 3) = zd.from_zip3)
      )
    WHERE dl.provides_services_counseling = true
      AND dl.id = p_duty_location_id
      AND j.code = p_gbloc_indicator
      AND toff.provides_ppm_closeout = true
    ORDER BY COALESCE(zd.distance_miles, 0) ASC;
END;
$$;

--B-22660 Daniel Jordan added fetch_counseling_offices_for_conus
CREATE OR REPLACE FUNCTION fetch_counseling_offices_for_conus(p_duty_location_id UUID)
RETURNS TABLE (id UUID, name TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        toff.id,
        toff.name
    FROM postal_code_to_gblocs pcg
    JOIN addresses a ON pcg.postal_code = a.postal_code
    JOIN duty_locations dl ON a.id = dl.address_id
    JOIN transportation_offices toff ON pcg.gbloc = toff.gbloc
    JOIN addresses toff_addr ON toff.address_id = toff_addr.id
    JOIN duty_locations d2 ON toff.id = d2.transportation_office_id
    JOIN re_us_post_regions rup ON toff_addr.postal_code = rup.uspr_zip_id
    LEFT JOIN zip3_distances zd
        ON (
            (substring(a.postal_code, 1, 3) = zd.from_zip3 AND substring(toff_addr.postal_code, 1, 3) = zd.to_zip3)
            OR
            (substring(a.postal_code, 1, 3) = zd.to_zip3 AND substring(toff_addr.postal_code, 1, 3) = zd.from_zip3)
        )
    WHERE dl.provides_services_counseling = true
      AND dl.id = p_duty_location_id
      AND d2.provides_services_counseling = true
    GROUP BY toff.id, toff.name, zd.distance_miles
    ORDER BY COALESCE(zd.distance_miles, 0), toff.name ASC;
END;
$$;

--B-22660 Daniel Jordan added get_counseling_offices
CREATE OR REPLACE FUNCTION get_counseling_offices(
    p_duty_location_id UUID,
    p_service_member_id UUID
)
RETURNS TABLE (id UUID, name TEXT)
LANGUAGE plpgsql AS $$
DECLARE
    is_address_oconus BOOLEAN;
    duty_address_id UUID;
    service_affiliation TEXT;
    dept_indicator TEXT;
    gbloc_indicator TEXT;
BEGIN

    SELECT duty_addr_id, is_oconus INTO duty_address_id, is_address_oconus
    FROM get_duty_location_info(p_duty_location_id);

    IF duty_address_id IS NULL THEN
        RAISE EXCEPTION 'Duty location % not found when searching for counseling offices', p_duty_location_id;
    END IF;

    IF is_address_oconus THEN
        service_affiliation := get_service_affiliation(p_service_member_id);
        dept_indicator := get_department_indicator(service_affiliation);

        gbloc_indicator := get_gbloc_indicator(duty_address_id, dept_indicator);

        RETURN QUERY SELECT * FROM fetch_counseling_offices_for_oconus(p_duty_location_id, gbloc_indicator);
    ELSE
        RETURN QUERY SELECT * FROM fetch_counseling_offices_for_conus(p_duty_location_id);
    END IF;
END;
$$;
