-- B-22468 M.Inthavongsay function to get the rate area for any address
CREATE OR REPLACE FUNCTION get_rate_area(
    address_id UUID,
    service_item_id UUID,
    c_id uuid
)
RETURNS TABLE (
    id UUID,
    contract_id UUID,
    is_oconus bool,
    code varchar(80),
    name varchar(80),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
) AS $$
DECLARE
    rate_area_id UUID;
BEGIN
    rate_area_id := get_rate_area_id(address_id, service_item_id, c_id);
    IF rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate Area ID not found for Address ID: % , Service Item ID: %, Contract ID: %', address_id, service_item_id, c_id;
    END IF;

    return query select re_rate_areas.id, re_rate_areas.contract_id, re_rate_areas.is_oconus,
        re_rate_areas.code, re_rate_areas.name, re_rate_areas.created_at, re_rate_areas.updated_at
        from re_rate_areas where re_rate_areas.id = rate_area_id;
END;
$$ LANGUAGE plpgsql;