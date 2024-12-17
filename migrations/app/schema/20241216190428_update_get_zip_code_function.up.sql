-- removing the exception that was previously being returned
-- this is to avoid the db update failing for the entire proc it is used in
-- we won't always have the POE/POD locations and want to ignore any errors here
CREATE OR REPLACE FUNCTION get_zip_code(shipment_id uuid, address_type VARCHAR)
RETURNS VARCHAR AS $$
    DECLARE zip_code VARCHAR;
    BEGIN

    IF address_type = 'pickup' THEN
        SELECT vl.uspr_zip_id
        INTO zip_code
        FROM mto_shipments ms
        JOIN addresses a ON a.id = ms.pickup_address_id
        JOIN v_locations vl ON vl.uprc_id = a.us_post_region_cities_id
        WHERE ms.id = shipment_id;
    ELSIF address_type = 'destination' THEN
        SELECT vl.uspr_zip_id
        INTO zip_code
        FROM mto_shipments ms
        JOIN addresses a ON a.id = ms.destination_address_id
        JOIN v_locations vl ON vl.uprc_id = a.us_post_region_cities_id
        WHERE ms.id = shipment_id;
    ELSIF address_type = 'poe' THEN
        SELECT vl.uspr_zip_id
        INTO zip_code
        FROM mto_service_items si
        JOIN port_locations pl ON pl.id = si.poe_location_id
        JOIN v_locations vl ON vl.uprc_id = pl.us_post_region_cities_id
        WHERE si.mto_shipment_id = shipment_id;
    ELSIF address_type = 'pod' THEN
        SELECT vl.uspr_zip_id
        INTO zip_code
        FROM mto_service_items si
        JOIN port_locations pl ON pl.id = si.pod_location_id
        JOIN v_locations vl ON vl.uprc_id = pl.us_post_region_cities_id
        WHERE si.mto_shipment_id = shipment_id;
    END IF;

    RETURN zip_code;
END;
$$ LANGUAGE plpgsql;