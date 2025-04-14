--B-22660 Daniel Jordan added get_service_affiliation
--B-22914 Beth Grohmann moved to separate script
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