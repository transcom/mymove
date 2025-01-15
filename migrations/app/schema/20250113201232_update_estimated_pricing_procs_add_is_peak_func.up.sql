-- removing ServiceAreaOrigin
DELETE FROM service_params
WHERE service_id = '9f3d551a-0725-430e-897e-80ee9add3ae9'
AND service_item_param_key_id = '599bbc21-8d1d-4039-9a89-ff52e3582144';

-- function that evaluates a date and returns T/F if it is during peak period
CREATE OR REPLACE FUNCTION is_peak_period(input_date DATE) RETURNS BOOLEAN AS $$
DECLARE
    peak_start DATE := MAKE_DATE(EXTRACT(YEAR FROM input_date)::INT, 5, 15); -- May 15th of the input year
    peak_end DATE := MAKE_DATE(EXTRACT(YEAR FROM input_date)::INT, 9, 30);   -- September 30th of the input year
BEGIN
    IF input_date IS NULL THEN
        RAISE EXCEPTION 'Input date cannot be NULL';
    END IF;
    -- if the input date is between May 15 and September 30 (inclusive), return true
    IF input_date BETWEEN peak_start AND peak_end THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;


-- adding the is_peak_period check to refine the price query further
CREATE OR REPLACE FUNCTION calculate_escalated_price(
    o_rate_area_id UUID,
    d_rate_area_id UUID,
    re_service_id UUID,
    c_id UUID,
    service_code TEXT,
    requested_pickup_date DATE
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
    peak_period BOOLEAN;
BEGIN
    -- we need to query the appropriate table based on the service code
    -- need to establish if the shipment is being moved during peak period
    peak_period := is_peak_period(requested_pickup_date);
    IF service_code IN ('ISLH', 'UBP') THEN
        SELECT rip.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_prices rip
        WHERE rip.origin_rate_area_id = o_rate_area_id AND rip.destination_rate_area_id = d_rate_area_id
          AND rip.service_id = re_service_id
          AND rip.contract_id = c_id
          AND rip.is_peak_period = peak_period;
    ELSE
        SELECT riop.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_other_prices riop
        WHERE (riop.rate_area_id = o_rate_area_id OR riop.rate_area_id = d_rate_area_id OR
            (o_rate_area_id IS NULL AND d_rate_area_id IS NULL))
        AND riop.service_id = re_service_id
        AND riop.contract_id = c_id
        AND riop.is_peak_period = peak_period;
    END IF;

    RAISE NOTICE '% per unit cents: %', service_code, per_unit_cents;
    IF per_unit_cents IS NULL THEN
        RAISE EXCEPTION 'No per unit cents found for service item id: %, origin rate area: %, dest rate area: %, and contract_id: %', re_service_id, o_rate_area_id, d_rate_area_id, c_id;
    END IF;

    SELECT rcy.escalation_compounded
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id
        AND requested_pickup_date BETWEEN rcy.start_date AND rcy.end_date;

    IF escalation_factor IS NULL THEN
        RAISE EXCEPTION 'Escalation factor not found for contract_id %', c_id;
    END IF;
    -- calculate the escalated price, return in dollars (dividing by 100)
    per_unit_cents := per_unit_cents / 100; -- putting in dollars
    escalated_price := ROUND(per_unit_cents * escalation_factor, 2); -- rounding to two decimals (100.00)

    RETURN escalated_price;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION IF EXISTS get_fuel_price(date);


-- updating get_fuel_price to return an INT instead of decimal, we were rounding too soon
CREATE OR REPLACE FUNCTION get_fuel_price(requested_pickup_date DATE)
RETURNS INTEGER AS $$
DECLARE
    fuel_price_in_cents INTEGER;
BEGIN

    SELECT fuel_price_in_millicents
    INTO fuel_price_in_cents
    FROM ghc_diesel_fuel_prices
    WHERE requested_pickup_date BETWEEN effective_date AND end_date;

    -- fallback to most recent fuel price if no match
    IF fuel_price_in_cents IS NULL THEN
        SELECT fuel_price_in_millicents
        INTO fuel_price_in_cents
        FROM ghc_diesel_fuel_prices
        ORDER BY publication_date DESC
        LIMIT 1;
    END IF;

    IF fuel_price_in_cents IS NULL THEN
        RAISE EXCEPTION 'No fuel price found for requested_pickup_date: %', requested_pickup_date;
    END IF;

    RAISE NOTICE 'Received fuel price of % for requested_pickup_date: %', fuel_price_in_cents, requested_pickup_date;

    RETURN fuel_price_in_cents;
END;
$$ LANGUAGE plpgsql;

-- updating to subtract the millicents value to avoid premature rounding
CREATE OR REPLACE FUNCTION calculate_price_difference(fuel_price DECIMAL)
RETURNS DECIMAL AS $$
BEGIN
    RETURN (fuel_price - 250000)::DECIMAL / 1000;
END;
$$ LANGUAGE plpgsql;

-- updating to use the shipment.requested_pickup_date value to refine search to get more accurate prices
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
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;

            WHEN service_code IN (''IHPK'', ''IUBPK'') THEN
                -- perform IHPK/IUBPK-specific logic (no destination rate area)
                contract_id := get_contract_id(shipment.requested_pickup_date);
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
			        -- update the pricing_estimate value in mto_service_items
			        UPDATE mto_service_items
			        SET pricing_estimate = estimated_price
			        WHERE id = service_item.id;
                END IF;

            WHEN service_code IN (''IHUPK'', ''IUBUPK'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                contract_id := get_contract_id(shipment.requested_pickup_date);
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id, contract_id);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, contract_id, service_code, shipment.requested_pickup_date);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND((escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC) * 100, 0);
                    RAISE NOTICE ''%: Received estimated price of % (% * (% / 100)) cents'', service_code, estimated_price, escalated_price, shipment.prime_estimated_weight;
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

                -- calculate estimated price, return as cents
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
        END CASE;
    END LOOP;
END;
'
LANGUAGE plpgsql;