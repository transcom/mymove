--B-22660 Daniel Jordan added get_gbloc_indicator
--B-22914 Beth Grohmann moved to separate script
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