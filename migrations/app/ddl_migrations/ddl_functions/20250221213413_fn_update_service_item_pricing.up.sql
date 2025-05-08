--B-22462  M.Inthavongsay Adding initial migration file for update_service_item_pricing stored procedure using new migration process.
--Also updating to allow IOSFSC and IDSFSC SIT service items.
--B-22742  C. Kleinjan  Add pricing calculation for ICRT and IUCRT service items
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
    distance NUMERIC;
    estimated_fsc_multiplier NUMERIC;
    fuel_price NUMERIC;
    cents_above_baseline NUMERIC;
    price_difference NUMERIC;
    days_in_sit INTEGER;
    length NUMERIC;
    width NUMERIC;
    height NUMERIC;
    standalone BOOLEAN;
    external BOOLEAN;
    cubic_feet NUMERIC;
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
        SELECT si.id, si.re_service_id, si.sit_delivery_miles,
        sit_origin_hhg_actual_address_id, sit_destination_final_address_id
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
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    -- multiply by 110% of estimated weight
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100), 2) * 100;
                    RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;

            WHEN service_code IN (''IHPK'', ''IUBPK'', ''IOSHUT'') THEN
                -- perform IHPK/IUBPK-specific logic (no destination rate area)
                contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    -- multiply by 110% of estimated weight
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100), 2) * 100;
                    RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;
            WHEN service_code IN (''ICRT'') THEN
                contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date, NULL);

                SELECT INTO length, height, width length_thousandth_inches, height_thousandth_inches, width_thousandth_inches
                FROM mto_service_item_dimensions
                WHERE mto_service_item_id = service_item.id AND type = ''CRATE'';

                SELECT INTO standalone, external standalone_crate, external_crate
                FROM mto_service_items
                WHERE id = service_item.id;

                IF length IS NOT NULL AND height IS NOT NULL AND width IS NOT NULL THEN
                    cubic_feet := ((length/1000) * (width/1000) * (height/1000)) / 1728;
                    
                    IF cubic_feet < 4.0 AND external THEN 
                        cubic_feet := 4.0;
                    END IF;

                    estimated_price := ROUND((escalated_price * cubic_feet), 2) * 100;

                    IF estimated_price > 100000.00 AND standalone THEN 
                        estimated_price := 100000.00;
                    END IF;

					RAISE NOTICE ''%: Received estimated price of % = (% * (%"L * %"W * %"H) / 100)) cents'', service_code, estimated_price, escalated_price, ROUND(length/1000, 3), ROUND(width/1000, 3), ROUND(height/1000, 3);
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;
            WHEN service_code IN (''IUCRT'') THEN
                contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date, NULL);

                SELECT INTO length, height, width length_thousandth_inches, height_thousandth_inches, width_thousandth_inches
                FROM mto_service_item_dimensions
                WHERE mto_service_item_id = service_item.id AND type = ''CRATE'';

                SELECT INTO standalone, external standalone_crate, external_crate
                FROM mto_service_items
                WHERE id = service_item.id;

                IF length IS NOT NULL AND height IS NOT NULL AND width IS NOT NULL THEN
                    cubic_feet := ((length/1000) * (width/1000) * (height/1000)) / 1728;
                    
                    IF cubic_feet < 4.0 AND external THEN 
                        cubic_feet := 4.0;
                    END IF;

                    estimated_price := ROUND((escalated_price * cubic_feet), 2) * 100;

                    IF estimated_price > 100000.00 AND standalone THEN 
                        estimated_price := 100000.00;
                    END IF;

					RAISE NOTICE ''%: Received estimated price of % = (% * (%"L * %"W * %"H) / 100)) cents'', service_code, estimated_price, escalated_price, ROUND(length/1000, 3), ROUND(width/1000, 3), ROUND(height/1000, 3);
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;
            WHEN service_code IN (''IHUPK'', ''IUBUPK'', ''IDSHUT'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    -- multiply by 110% of estimated weight
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100), 2) * 100;
                    RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;

            WHEN service_code IN (''POEFSC'', ''PODFSC'') THEN
                distance = mileage;

                -- getting FSC multiplier from re_fsc_multipliers
                estimated_fsc_multiplier := get_fsc_multiplier(shipment.prime_estimated_weight);

                fuel_price := get_fuel_price(shipment.requested_pickup_date);

                price_difference := calculate_price_difference(fuel_price);

                IF estimated_fsc_multiplier IS NOT NULL AND distance IS NOT NULL THEN
                    cents_above_baseline := distance * estimated_fsc_multiplier;
                    RAISE NOTICE ''Distance: % * FSC Multipler: % = $% cents above baseline of $2.50'', distance, estimated_fsc_multiplier, cents_above_baseline;
                    RAISE NOTICE ''The fuel price is % above the baseline (% - 250000 baseline)'', price_difference, fuel_price;
                    estimated_price := ROUND((cents_above_baseline * price_difference) * 100);
                    RAISE NOTICE ''Received estimated price of % cents for service_code: %.'', estimated_price, service_code;

			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;

            WHEN service_code IN (''IOSFSC'', ''IDSFSC'') THEN
                distance = service_item.sit_delivery_miles;

                -- Pricing will not be executed if origin pickup is OCONUS. This is achieved with ZERO mileage in the calculation.
                IF service_code = ''IOSFSC'' AND service_item.sit_origin_hhg_actual_address_id IS NOT NULL THEN
                    IF get_is_oconus(service_item.sit_origin_hhg_actual_address_id) THEN
                        distance := 0;
                        RAISE NOTICE ''Pickup[service_item.sit_origin_hhg_actual_address_id: %] is OCONUS. Distance will be set to 0 to cause pricing to be 0 cents'', service_item.sit_origin_hhg_actual_address_id;
                    END IF;
                END IF;

                -- Pricing will not be executed if origin destination is OCONUS. This is achieved with ZERO mileage in the calculation.
                IF service_code = ''IDSFSC'' AND service_item.sit_destination_final_address_id IS NOT NULL THEN
                    IF get_is_oconus(service_item.sit_destination_final_address_id) THEN
                        distance := 0;
                        RAISE NOTICE ''Destination[service_item.sit_destination_final_address_id: %] is OCONUS. Distance will be set to 0 to cause pricing to be 0 cents'', service_item.sit_destination_final_address_id;
                    END IF;
                END IF;

                estimated_fsc_multiplier := get_fsc_multiplier(shipment.prime_estimated_weight);

                fuel_price := get_fuel_price(shipment.requested_pickup_date);

                price_difference := calculate_price_difference(fuel_price);

                IF estimated_fsc_multiplier IS NOT NULL AND distance IS NOT NULL THEN
                    cents_above_baseline := distance * estimated_fsc_multiplier;
                    RAISE NOTICE ''Distance: % * FSC Multipler: % = $% cents above baseline of $2.50'', distance, estimated_fsc_multiplier, cents_above_baseline;
                    RAISE NOTICE ''The fuel price is % above the baseline (% - 250000 baseline)'', price_difference, fuel_price;
                    estimated_price := ROUND((cents_above_baseline * price_difference) * 100);
                    RAISE NOTICE ''Received estimated price of % cents for service_code: %.'', estimated_price, service_code;

			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                ELSE
                    RAISE NOTICE ''service_code: % - Failed to compute pricing[estimated_fsc_multiplier: %, distance: %]'', service_code, estimated_fsc_multiplier, distance;
                END IF;
            ELSE
                RAISE warning ''Unsupported service code: %'', service_code;
        END CASE;
    END LOOP;
END;
'
LANGUAGE plpgsql;