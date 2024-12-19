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


-- updating the get rate area function to include the contract id
CREATE OR REPLACE FUNCTION get_rate_area_id(
    address_id UUID,
    service_item_id UUID,
    c_id uuid,
    OUT o_rate_area_id UUID
)
RETURNS UUID AS $$
DECLARE
    is_oconus BOOLEAN;
    zip3_value TEXT;
BEGIN
    is_oconus := get_is_oconus(address_id);

    IF is_oconus THEN
        -- re_oconus_rate_areas if is_oconus is TRUE
        SELECT ro.rate_area_id
        INTO o_rate_area_id
        FROM addresses a
        JOIN re_oconus_rate_areas ro
        ON a.us_post_region_cities_id = ro.us_post_region_cities_id
        JOIN re_rate_areas ra ON ro.rate_area_id = ra.id
        WHERE a.id = address_id
            AND ra.contract_id = c_id;
    ELSE
        -- re_zip3s if is_oconus is FALSE
        SELECT rupr.zip3
        INTO zip3_value
        FROM addresses a
        JOIN us_post_region_cities uprc
        ON a.us_post_region_cities_id = uprc.id
        JOIN re_us_post_regions rupr
        ON uprc.us_post_regions_id = rupr.id
        WHERE a.id = address_id;

        -- use the zip3 value to find the rate_area_id in re_zip3s
        SELECT rz.rate_area_id
        INTO o_rate_area_id
        FROM re_zip3s rz
        JOIN re_rate_areas ra
        ON rz.rate_area_id = ra.id
        WHERE rz.zip3 = zip3_value
            AND ra.contract_id = c_id;
    END IF;

    -- Raise an exception if no rate area is found
    IF o_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate area not found for address % for service item ID %', address_id, service_item_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- this function will help us get the ZIP code & port type for a port to calculate mileage
CREATE OR REPLACE FUNCTION get_port_location_info_for_shipment(shipment_id UUID)
RETURNS TABLE(uspr_zip_id TEXT, port_type TEXT) AS $$
BEGIN
    -- select the ZIP code and port type (POEFSC or PODFSC)
    RETURN QUERY
    SELECT
        COALESCE(poe_usprc.uspr_zip_id::TEXT, pod_usprc.uspr_zip_id::TEXT) AS uspr_zip_id,
        CASE
            WHEN msi.poe_location_id IS NOT NULL THEN 'POEFSC'
            WHEN msi.pod_location_id IS NOT NULL THEN 'PODFSC'
            ELSE NULL
        END AS port_type
    FROM mto_shipments ms
    JOIN mto_service_items msi ON ms.id = msi.mto_shipment_id
    LEFT JOIN port_locations poe_pl ON msi.poe_location_id = poe_pl.id
    LEFT JOIN port_locations pod_pl ON msi.pod_location_id = pod_pl.id
    LEFT JOIN us_post_region_cities poe_usprc ON poe_pl.us_post_region_cities_id = poe_usprc.id
    LEFT JOIN us_post_region_cities pod_usprc ON pod_pl.us_post_region_cities_id = pod_usprc.id
    WHERE ms.id = shipment_id
    AND (msi.poe_location_id IS NOT NULL OR msi.pod_location_id IS NOT NULL)
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- updating the pricing proc to now consume the mileage we get from DTOD instead of calculate it using Rand McNally
-- this is a requirement for E-06210
-- also updating the get_rate_area parameters and passing in the contract_id
CREATE OR REPLACE PROCEDURE update_service_item_pricing(
    shipment_id UUID,
    mileage INT
) AS
'
DECLARE
    shipment RECORD;
    service_item RECORD;
    escalated_price NUMERIC;
    estimated_price NUMERIC;
    o_rate_area_id UUID;
    d_rate_area_id UUID;
    contract_id UUID;
    service_code TEXT;
    o_zip_code TEXT;
    d_zip_code TEXT;
    distance NUMERIC;  -- This will be replaced by mileage
    estimated_fsc_multiplier NUMERIC;
    fuel_price NUMERIC;
    cents_above_baseline NUMERIC;
    price_difference NUMERIC;
BEGIN
    SELECT ms.id, ms.pickup_address_id, ms.destination_address_id, ms.requested_pickup_date, ms.prime_estimated_weight
    INTO shipment
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF shipment IS NULL THEN
        RAISE EXCEPTION ''Shipment with ID % not found'', shipment_id;
    END IF;

    -- exit the proc if prime_estimated_weight is NULL
    IF shipment.prime_estimated_weight IS NULL THEN
        RETURN;
    END IF;

    -- loop through service items in the shipment
    FOR service_item IN
        SELECT si.id, si.re_service_id
        FROM mto_service_items si
        WHERE si.mto_shipment_id = shipment_id
    LOOP
        -- get the service code for the current service item to determine calculation
        SELECT code
        INTO service_code
        FROM re_services
        WHERE id = service_item.re_service_id;

        CASE
            WHEN service_code IN (''ISLH'', ''UBP'') THEN
                contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, contract_id);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, contract_id, service_code);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
                END IF;

            WHEN service_code IN (''IHPK'', ''IUBPK'') THEN
                -- perform IHPK/IUBPK-specific logic (no destination rate area)
                contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, contract_id, service_code);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
                END IF;

            WHEN service_code IN (''IHUPK'', ''IUBUPK'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, contract_id, service_code);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
                END IF;

            WHEN service_code IN (''POEFSC'', ''PODFSC'') THEN
                -- use the passed mileage parameter
                distance := mileage;

                -- getting FSC multiplier from re_fsc_multipliers
                estimated_fsc_multiplier := get_fsc_multiplier(shipment.prime_estimated_weight);

                fuel_price := get_fuel_price(shipment.requested_pickup_date);

                price_difference := calculate_price_difference(fuel_price);

                -- calculate estimated price, return as cents
                IF estimated_fsc_multiplier IS NOT NULL AND distance IS NOT NULL THEN
                    cents_above_baseline := distance * estimated_fsc_multiplier;
                    RAISE NOTICE ''Distance: % * FSC Multipler: % = $% cents above baseline of $2.50'', distance, estimated_fsc_multiplier, cents_above_baseline;
                    RAISE NOTICE ''The fuel price is % cents above the baseline ($% - $2.50 baseline)'', price_difference, fuel_price;
                    estimated_price := ROUND((cents_above_baseline * price_difference) * 100);
                    RAISE NOTICE ''Received estimated price of % cents for service_code: %.'', estimated_price, service_code;
                END IF;
        END CASE;

        -- update the pricing_estimate value in mto_service_items
        UPDATE mto_service_items
        SET pricing_estimate = estimated_price
        WHERE id = service_item.id;
    END LOOP;
END;
'
LANGUAGE plpgsql;
