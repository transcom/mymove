--B-22660 Daniel Jordan added get_counseling_offices
--B-22914 Beth Grohmann moved other functions to separate scripts
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
