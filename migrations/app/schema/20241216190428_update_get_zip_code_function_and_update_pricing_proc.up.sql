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

-- updating the pricing proc to now consume the mileage we get from DTOD instead of calculate it using Rand McNally
-- this is a requirement for E-06210
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
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id);
                contract_id := get_contract_id(shipment.requested_pickup_date);
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, contract_id, service_code);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
                END IF;

            WHEN service_code IN (''IHPK'', ''IUBPK'') THEN
                -- perform IHPK/IUBPK-specific logic (no destination rate area)
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id);
                contract_id := get_contract_id(shipment.requested_pickup_date);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, contract_id, service_code);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
                END IF;

            WHEN service_code IN (''IHUPK'', ''IUBUPK'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id);
                contract_id := get_contract_id(shipment.requested_pickup_date);
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
