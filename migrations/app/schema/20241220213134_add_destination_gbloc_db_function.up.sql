-- this function will handle getting the destination GBLOC associated with a shipment's destination address
-- this only applies to OCONUS destination addresses on a shipment, but this can also checks domestic shipments
CREATE OR REPLACE FUNCTION get_destination_gbloc_for_shipment(shipment_id UUID)
RETURNS TEXT AS $$
DECLARE
    service_member_affiliation TEXT;
    zip TEXT;
    gbloc_result TEXT;
    alaska_zone_ii BOOLEAN;
    market_code TEXT;
BEGIN
    -- get the shipment's market code to determine conditionals
    SELECT ms.market_code
    INTO market_code
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    -- if it's a domestic shipment, use postal_code_to_gblocs
    IF market_code = 'd' THEN
        SELECT upc.uspr_zip_id
        INTO zip
        FROM addresses a
        JOIN us_post_region_cities upc ON a.us_post_region_cities_id = upc.id
        WHERE a.id = (SELECT destination_address_id FROM mto_shipments WHERE id = shipment_id);

        SELECT gbloc
        INTO gbloc_result
        FROM postal_code_to_gblocs
        WHERE postal_code = zip
        LIMIT 1;

        IF gbloc_result IS NULL THEN
            RETURN NULL;
        END IF;

        RETURN gbloc_result;

    ELSEIF market_code = 'i' THEN
        -- if it's 'i' then we need to check for some exceptions
        SELECT sm.affiliation
        INTO service_member_affiliation
        FROM service_members sm
        JOIN orders o ON o.service_member_id = sm.id
        JOIN moves m ON m.orders_id = o.id
        JOIN mto_shipments ms ON ms.move_id = m.id
        WHERE ms.id = shipment_id;

        SELECT upc.uspr_zip_id
        INTO zip
        FROM addresses a
        JOIN us_post_region_cities upc ON a.us_post_region_cities_id = upc.id
        WHERE a.id = (SELECT destination_address_id FROM mto_shipments WHERE id = shipment_id);

        -- check if the postal code (uspr_zip_id) is in Alaska Zone II
        SELECT EXISTS (
            SELECT 1
            FROM re_oconus_rate_areas ro
            JOIN re_rate_areas ra ON ro.rate_area_id = ra.id
            JOIN us_post_region_cities upc ON upc.id = ro.us_post_region_cities_id
            WHERE upc.uspr_zip_id = zip
              AND ra.code = 'US8190100'  -- Alaska Zone II Code
        )
        INTO alaska_zone_ii;

        -- if the service member is USAF or USSF and the address is in Alaska Zone II, return 'MBFL'
        IF (service_member_affiliation = 'AIR_FORCE' OR service_member_affiliation = 'SPACE_FORCE') AND alaska_zone_ii THEN
            RETURN 'MBFL';
        END IF;

        -- for all other branches except USMC, return the gbloc from the postal_code_to_gbloc table based on the zip
        SELECT gbloc
        INTO gbloc_result
        FROM postal_code_to_gblocs
        WHERE postal_code = zip
        LIMIT 1;

        IF gbloc_result IS NULL THEN
            RETURN NULL;
        END IF;

        RETURN gbloc_result;
    END IF;

END;
$$ LANGUAGE plpgsql;
