
--B-22462  M.Inthavongsay Adding initial migration file for update_service_item_pricing stored procedure using new migration process.
--Also updating to allow IOSFSC and IDSFSC SIT service items.
--B-22463  M.Inthavongsay updating to allow IOASIT and IDASIT SIT service items.
--B-22662  C.Jewell added INPK estimate pricing
--B-22466  M.Inthavongsay updating to allow IOPSIT and IDDSIT SIT service items.
--B-22464  A Lusk updating to allow IOFSIT and IDFSIT service items.
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
    service_code TEXT;
    o_zip_code TEXT;
    d_zip_code TEXT;
    distance NUMERIC;
    estimated_fsc_multiplier NUMERIC;
    fuel_price NUMERIC;
    cents_above_baseline NUMERIC;
    price_difference NUMERIC;
    days_in_sit INTEGER;
    declared_contract_id UUID;
    declared_escalation_factor NUMERIC;
    declared_oconus_factor NUMERIC;
    declared_market_code TEXT;
    declared_is_oconus BOOLEAN;
    length NUMERIC;
    width NUMERIC;
    height NUMERIC;
    standalone BOOLEAN;
    external BOOLEAN;
    cubic_feet NUMERIC;
    standalone_crate_cap NUMERIC;
    external_crate_minimum NUMERIC;
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

    SELECT parameter_value
    INTO standalone_crate_cap
    FROM application_parameters
    WHERE parameter_name = ''standaloneCrateCap'';

    SELECT parameter_value
    INTO external_crate_minimum
    FROM application_parameters
    WHERE parameter_name = ''externalCrateMinimum'';

    -- loop through service items in the shipment
    FOR service_item IN
        SELECT si.id, si.re_service_id, si.sit_delivery_miles, si.sit_departure_date, si.sit_entry_date,
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
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);

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
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);

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
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);

                SELECT INTO length, height, width length_thousandth_inches, height_thousandth_inches, width_thousandth_inches
                FROM mto_service_item_dimensions
                WHERE mto_service_item_id = service_item.id AND type = ''CRATE'';

                SELECT INTO standalone, external standalone_crate, external_crate
                FROM mto_service_items
                WHERE id = service_item.id;

                IF length IS NOT NULL AND height IS NOT NULL AND width IS NOT NULL THEN
                    cubic_feet := ROUND(((length/1000) * (width/1000) * (height/1000)) / 1728, 2);

                    IF cubic_feet < external_crate_minimum AND external THEN
                        cubic_feet := external_crate_minimum;
                    END IF;

                    estimated_price := ROUND((escalated_price * cubic_feet), 2) * 100;

                    IF estimated_price > standalone_crate_cap AND standalone THEN
                        estimated_price := standalone_crate_cap;
                    END IF;

					RAISE NOTICE ''%: Received estimated price of % cents = %¢/ft³  * %ft³'', service_code, estimated_price, escalated_price, cubic_feet;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;
            WHEN service_code IN (''IUCRT'') THEN
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(d_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);

                SELECT INTO length, height, width length_thousandth_inches, height_thousandth_inches, width_thousandth_inches
                FROM mto_service_item_dimensions
                WHERE mto_service_item_id = service_item.id AND type = ''CRATE'';

                SELECT INTO standalone, external standalone_crate, external_crate
                FROM mto_service_items
                WHERE id = service_item.id;

                IF length IS NOT NULL AND height IS NOT NULL AND width IS NOT NULL THEN
                    cubic_feet := ROUND(((length/1000) * (width/1000) * (height/1000)) / 1728, 2);

                    IF cubic_feet < external_crate_minimum AND external THEN
                        cubic_feet := external_crate_minimum;
                    END IF;

                    estimated_price := ROUND((escalated_price * cubic_feet), 2) * 100;

                    IF estimated_price > standalone_crate_cap AND standalone THEN
                        estimated_price := standalone_crate_cap;
                    END IF;

					RAISE NOTICE ''%: Received estimated price of % cents = %¢/ft³  * %ft³'', service_code, estimated_price, escalated_price, cubic_feet;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;
            WHEN service_code IN (''IHUPK'', ''IUBUPK'', ''IDSHUT'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);

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

            WHEN service_code IN (''IOASIT'', ''IDASIT'') THEN
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);

                IF service_code = ''IOASIT'' THEN
                    o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                    escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);
                ELSE
                    d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, declared_contract_id);
                    escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);
                END IF;

                BEGIN
                    -- Retrieve MAX days in sit allowance value from application parameter table.
                    days_in_sit := get_application_parameter_value(''maxSitDaysAllowance'')::int - 1;
                EXCEPTION WHEN OTHERS THEN
                    RAISE EXCEPTION ''%: unexpected error parsing maxSitDaysAllowance application param value'', service_code;
                END;

                IF days_in_sit IS NULL THEN
                    RAISE EXCEPTION ''%: maxSitDaysAllowance application param value not found'', service_code;
                END IF;

                IF service_item.sit_entry_date IS NOT NULL AND service_item.sit_departure_date IS NOT NULL THEN
                    days_in_sit := (SELECT (service_item.sit_departure_date::date - (service_item.sit_entry_date::date)) as days);
                END IF;

                RAISE NOTICE ''days_in_sit = %'', days_in_sit;

                IF escalated_price IS NOT NULL AND days_in_sit IS NOT NULL AND days_in_sit >= 0 THEN
                    RAISE NOTICE ''escalated_price = $% cents'', escalated_price;

                    -- multiply by 110% of estimated weight
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100) * days_in_sit, 2) * 100;
                    RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100) * %) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight, days_in_sit;

                    -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                ELSE
                    RAISE NOTICE ''service_code: % - Failed to compute pricing[escalated_price: %, days_in_sit: %]'', service_code, escalated_price, days_in_sit;
                END IF;
            WHEN service_code = ''INPK'' THEN
                -- INPK requires the base price for an origin rate area and a requested pickup date
                -- get the base price for the origin rate area from IHPK (iHHG into iNTS means use IHPK base price)
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);

                -- Use IHPK for the escalated price for the INPK case
                -- This is because the scenario is iHHG -> iNTS
                escalated_price := calculate_escalated_price(
                    o_rate_area_id,
                    NULL,
                    (SELECT id FROM re_services WHERE code = ''IHPK''),
                    declared_contract_id,
                    ''IHPK'',
                    shipment.requested_pickup_date,
                    NULL
                );

                -- Now that we have the escalated price, we multiply it by the
                -- NTS INPK market code factor. This time we pass in INPK,
                -- because this is an NTS scenario
                declared_oconus_factor := get_market_code_factor_escalation(
                    o_rate_area_id,
                    declared_contract_id,
                    service_item.re_service_id
                );

                -- Okay, now that we have all of our numbers. We just gotta calc
                -- the final price

                -- Final estimated price = escalated price * factor * 110% of estimated weight
                estimated_price := ROUND(
                    ( escalated_price * declared_oconus_factor * ((shipment.prime_estimated_weight * 1.1) / 100.0) )::numeric,
                    2
                ) * 100;

                RAISE NOTICE ''INPK: esc=%, factor=%, cwt=%, final=% (service_item id=%)'',
                    escalated_price,
                    declared_oconus_factor,
                    (shipment.prime_estimated_weight / 100.0),
                    estimated_price,
                    service_item.id;

                UPDATE mto_service_items
                SET pricing_estimate = estimated_price
                WHERE id = service_item.id;
            WHEN service_code IN (''IOPSIT'', ''IDDSIT'') THEN
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);

                distance = service_item.sit_delivery_miles;
                RAISE NOTICE ''SIT mileage = %'', distance;

                IF service_code = ''IOPSIT'' THEN
                    o_rate_area_id := get_rate_area_id(service_item.sit_origin_hhg_actual_address_id, service_item.re_service_id, declared_contract_id);
                    escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, distance);
                ELSE
                    d_rate_area_id := get_rate_area_id(service_item.sit_destination_final_address_id, service_item.re_service_id, declared_contract_id);
                    escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, distance);
                END IF;

                IF distance IS NOT NULL AND escalated_price IS NOT NULL AND shipment.prime_estimated_weight IS NOT NULL THEN
                    RAISE NOTICE ''escalated_price = $% cents'', escalated_price;

                    IF distance > 50 THEN
                        -- multiply by 110% of estimated weight
                        -- multiply by mileage
                        estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100) * distance, 2) * 100;

                        RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100) * %) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight, distance;
                    ELSE
                        -- multiply by 110% of estimated weight
                        estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100), 2) * 100;

                        RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
                    END IF;

                    -- update the pricing_estimate value in mto_service_items
                    UPDATE mto_service_items
                    SET pricing_estimate = estimated_price
                    WHERE id = service_item.id;
                ELSE
                    RAISE NOTICE ''service_code: % - Failed to compute pricing[escalated_price: %, prime_estimated_weight: %, distance: %]'', service_code, escalated_price, shipment.prime_estimated_weight, distance;
                END IF;

			WHEN service_code IN (''IOFSIT'', ''IDFSIT'') THEN
				declared_contract_id := get_contract_id(shipment.requested_pickup_date);

				IF service_code = ''IOFSIT'' THEN
                    o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                    escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);
                ELSE
                    d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, declared_contract_id);
                    escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date, NULL);
                END IF;

				IF escalated_price IS NOT NULL THEN
                    RAISE NOTICE ''escalated_price = $% cents'', escalated_price;

                    -- multiply by 110% of estimated weight
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100), 2) * 100;
                    RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;

                    -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                ELSE
                    RAISE NOTICE ''service_code: % - Failed to compute pricing[escalated_price: %, days_in_sit: %]'', service_code, escalated_price, days_in_sit;
                END IF;

            ELSE
                RAISE warning ''Unsupported service code: %'', service_code;
        END CASE;
    END LOOP;
END;
'
LANGUAGE plpgsql;
