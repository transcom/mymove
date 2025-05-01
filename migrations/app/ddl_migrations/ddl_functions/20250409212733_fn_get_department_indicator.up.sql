--B-22660 Daniel Jordan added get_department_indicator
--B-22914 Beth Grohmann moved to separate script
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