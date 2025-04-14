--B-22660 Daniel Jordan added get_duty_location_info
--B-22914 Beth Grohmann moved to separate script
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