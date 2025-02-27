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
    declared_contract_id UUID;
    service_code TEXT;
    o_zip_code TEXT;
    d_zip_code TEXT;
    distance NUMERIC;
    estimated_fsc_multiplier NUMERIC;
    fuel_price NUMERIC;
    cents_above_baseline NUMERIC;
    price_difference NUMERIC;
    declared_base_price NUMERIC;
    declared_escalation_factor NUMERIC;
    declared_oconus_factor NUMERIC;
    declared_market_code TEXT;
    declared_is_oconus BOOLEAN;
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
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date);

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
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    -- multiply by 110% of estimated weight
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight * 1.1) / 100), 2) * 100;
                    RAISE NOTICE ''%: Received estimated price of % (% * (% * 1.1) / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;

            WHEN service_code IN (''IHUPK'', ''IUBUPK'', ''IDSHUT'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, declared_contract_id);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, declared_contract_id, service_code, shipment.requested_pickup_date);

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
            WHEN service_code = ''INPK'' THEN
                -- INPK requires the base price for an origin rate area and a requested pickup date
                -- get the base price for the origin rate area from IHPK (iHHG into iNTS means use IHPK base price)
                declared_contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, declared_contract_id);
                SELECT riop.per_unit_cents
                INTO declared_base_price
                FROM re_intl_other_prices AS riop
                JOIN re_contract_years AS rcy ON rcy.contract_id = riop.contract_id
                WHERE riop.contract_id = declared_contract_id
                AND riop.service_id = (SELECT id FROM re_services WHERE code = ''IHPK'' LIMIT 1)
                AND riop.rate_area_id = o_rate_area_id
                AND shipment.requested_pickup_date >= rcy.start_date
                AND shipment.requested_pickup_date <= rcy.end_date
                LIMIT 1;

                IF declared_base_price IS NULL THEN
                    RAISE EXCEPTION ''No base price found for IHPK when calculating INPK estimate price: o_rate_area_id=%, declared_contract_id=%, cwt=%, shipment.requested_pickup_date=% (service_item id=%)'',
                    o_rate_area_id,
                    declared_contract_id,
                    (shipment.prime_estimated_weight / 100.0),
                    shipment.requested_pickup_date,
                    service_item.id;
                    CONTINUE;
                END IF;
                -- Now that we have the IHPK base price
                -- we can get the escalation factor and thus the escalated price
                declared_escalation_factor := calculate_escalation_factor(declared_contract_id, shipment.requested_pickup_date);
                escalated_price := declared_base_price * declared_escalation_factor;

                -- Now that we have the escalated price, we multiply it by the
                -- NTS INPK market code factor

                SELECT is_oconus
                INTO declared_is_oconus
                FROM re_rate_areas
                WHERE id = o_rate_area_id;

                IF declared_is_oconus THEN
                    declared_market_code := ''O'';
                ELSE
                    declared_market_code := ''C'';
                END IF;

                SELECT stp.factor
                    INTO declared_oconus_factor
                    FROM re_shipment_type_prices stp
                 WHERE stp.contract_id = declared_contract_id
                 -- Use INPK for this one, not IHPK as we are applying
                 -- NTS math to IHPK price
                    AND stp.service_id = service_item.re_service_id
                    AND stp.market = declared_market_code
                    LIMIT 1;

                IF declared_oconus_factor IS NULL THEN
                    RAISE EXCEPTION ''No OCONUS/CONUS factor found for INPK for market_code=%, contract_id=%, re_service_id=%.'', declared_market_code, declared_contract_id, service_item.re_service_id;
                    CONTINUE;
                END IF;

                -- Okay, now that we have all of our numbers. We just gotta calc
                -- the final price

                -- Final estimated price = escalated price * factor * 110% of estimated weight
                estimated_price := ROUND(
                    ( escalated_price * declared_oconus_factor * ((shipment.prime_estimated_weight * 1.1) / 100.0) )::numeric,
                    2
                ) * 100;

                RAISE NOTICE ''INPK: base=%, esc=%, factor=%, cwt=%, final=% (service_item id=%)'',
                    declared_base_price,
                    escalated_price,
                    declared_oconus_factor,
                    (shipment.prime_estimated_weight / 100.0),
                    estimated_price,
                    service_item.id;

                UPDATE mto_service_items
                   SET pricing_estimate = estimated_price
                 WHERE id = service_item.id;
            ELSE
                -- Case if none of the above are triggered
                RAISE EXCEPTION ''Unknown service code: % for service item %'', service_code, service_item.id;
        END CASE;
    END LOOP;
END;
'
LANGUAGE plpgsql;