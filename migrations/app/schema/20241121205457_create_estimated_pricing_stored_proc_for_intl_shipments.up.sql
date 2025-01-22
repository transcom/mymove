-- function to get the rate area id for any address
CREATE OR REPLACE FUNCTION get_rate_area_id(
    address_id UUID,
    service_item_id UUID,
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
        WHERE a.id = address_id;
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
        WHERE rz.zip3 = zip3_value;
    END IF;

    -- Raise an exception if no rate area is found
    IF o_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate area not found for address % for service item ID %', address_id, service_item_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- function to get the contract id based off of a specific date that falls between start/end dates
CREATE OR REPLACE FUNCTION get_contract_id(
    requested_pickup_date DATE,
    OUT o_contract_id UUID
)
RETURNS UUID AS $$
BEGIN
    -- get the contract_id from the re_contract_years table
    SELECT rcy.contract_id
    INTO o_contract_id
    FROM re_contract_years rcy
    WHERE requested_pickup_date BETWEEN rcy.start_date AND rcy.end_date;

    -- check if contract_id is found, else raise an exception
    IF o_contract_id IS NULL THEN
        RAISE EXCEPTION 'Contract not found for requested pickup date %', requested_pickup_date;
    END IF;
END;
$$ LANGUAGE plpgsql;


-- function to calculate the escalated price, takes in:
-- origin rate area
-- dest rate area
-- re_services id
-- contract id
CREATE OR REPLACE FUNCTION calculate_escalated_price(
    o_rate_area_id UUID,
    d_rate_area_id UUID,
    re_service_id UUID,
    c_id UUID,
    service_code TEXT
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
BEGIN
    -- we need to query the appropriate table based on the service code
    IF service_code IN ('ISLH', 'UBP') THEN
        SELECT rip.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_prices rip
        WHERE (rip.origin_rate_area_id = o_rate_area_id OR o_rate_area_id IS NULL)
          AND (rip.destination_rate_area_id = d_rate_area_id OR d_rate_area_id IS NULL)
          AND rip.service_id = re_service_id
          AND rip.contract_id = c_id;
    ELSE
        SELECT riop.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_other_prices riop
        WHERE (riop.rate_area_id = o_rate_area_id OR riop.rate_area_id = d_rate_area_id OR
            (o_rate_area_id IS NULL AND d_rate_area_id IS NULL))
        AND riop.service_id = re_service_id
        AND riop.contract_id = c_id;

    END IF;

    IF per_unit_cents IS NULL THEN
        RAISE EXCEPTION 'No per unit cents found for service item id: %, origin rate area: %, dest rate area: %, and contract_id: %', re_service_id, o_rate_area_id, d_rate_area_id, c_id;
    END IF;

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

    IF escalation_factor IS NULL THEN
        RAISE EXCEPTION 'Escalation factor not found for contract_id %', c_id;
    END IF;
    -- calculate the escalated price, return in dollars (dividing by 100)
    escalated_price := ROUND(per_unit_cents * escalation_factor::NUMERIC / 100, 2);

    RAISE NOTICE '% escalated price: $% (% * % / 100)', service_code, escalated_price, per_unit_cents, escalation_factor;

    RETURN escalated_price;
END;
$$ LANGUAGE plpgsql;


-- get ZIP code by passing in a shipment ID and the address type
-- used for PODFSC & POEFSC service item types
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

    IF zip_code IS NULL THEN
        RAISE EXCEPTION 'zip_code not found for shipment id: % and address type of: %', shipment_id, address_type;
    END IF;

    RETURN zip_code;
END;
$$ LANGUAGE plpgsql;

-- getting the distance between two ZIPs
CREATE OR REPLACE FUNCTION get_distance(o_zip_code VARCHAR, d_zip_code VARCHAR)
RETURNS INT AS $$
DECLARE
    dist INT;
BEGIN
    -- get the last 3 characters from both zip codes
    SELECT zd.distance_miles
    INTO dist
    FROM zip3_distances zd
    WHERE (zd.from_zip3 = LEFT(o_zip_code, 3) AND zd.to_zip3 = LEFT(d_zip_code, 3))
       OR (zd.from_zip3 = LEFT(d_zip_code, 3) AND zd.to_zip3 = LEFT(o_zip_code, 3));

    IF dist IS NOT NULL THEN
        RAISE NOTICE 'Distance found between o_zip_code % and d_zip_code %: % miles', o_zip_code, d_zip_code, dist;
    ELSE
        RAISE NOTICE 'No distance found for o_zip_code: % and d_zip_code: %', o_zip_code, d_zip_code;
        RETURN 0;
    END IF;

    -- If no distance found, return 0 - we will have the backend call DTOD
    IF dist IS NULL THEN
        RETURN 0;
    END IF;

    RETURN dist;
END;
$$ LANGUAGE plpgsql;


-- querying the re_fsc_multiplier table and getting the multiplier value
CREATE OR REPLACE FUNCTION get_fsc_multiplier(estimated_weight INT)
RETURNS DECIMAL AS $$
    DECLARE m NUMERIC;
    BEGIN

    SELECT multiplier
    INTO m
    FROM re_fsc_multipliers
    WHERE estimated_weight >= low_weight AND estimated_weight <= high_weight;

    IF m IS NULL THEN
        RAISE EXCEPTION 'multipler not found for weight of %', estimated_weight;
    END IF;

    RAISE NOTICE 'Received FSC multiplier for estimated_weight: %', m;

    RETURN m;
END;
$$ LANGUAGE plpgsql;


-- getting the fuel price from the ghc_diesel_fuel_prices table
CREATE OR REPLACE FUNCTION get_fuel_price(requested_pickup_date DATE)
RETURNS DECIMAL AS $$
DECLARE
    fuel_price DECIMAL;
BEGIN

    SELECT ROUND(fuel_price_in_millicents::DECIMAL / 100000, 2)
    INTO fuel_price
    FROM ghc_diesel_fuel_prices
    WHERE requested_pickup_date BETWEEN effective_date AND end_date;

    -- if no results, fallback to the most recent fuel price
    IF fuel_price IS NULL THEN
        SELECT ROUND(fuel_price_in_millicents::DECIMAL / 100000, 2)
        INTO fuel_price
        FROM ghc_diesel_fuel_prices
        ORDER BY publication_date DESC
        LIMIT 1;
    END IF;

    IF fuel_price IS NULL THEN
        RAISE EXCEPTION 'No fuel price found for requested_pickup_date: %', requested_pickup_date;
    END IF;

    RAISE NOTICE 'Received fuel price of $% for requested_pickup_date: %', fuel_price, requested_pickup_date;

    RETURN fuel_price;
END;
$$ LANGUAGE plpgsql;


-- calculating difference from fuel price from base price, return in cents
CREATE OR REPLACE FUNCTION calculate_price_difference(fuel_price DECIMAL)
RETURNS DECIMAL AS $$
BEGIN
    RETURN (fuel_price - 2.50) * 100;
END;
$$ LANGUAGE plpgsql;

-- function that handles calculating price for ISLH & UBP service items, takes in:
-- origin rate area
-- dest rate area
-- re_services id
-- contract id
-- prime estimated weight
CREATE OR REPLACE FUNCTION calculate_islh_ubp_price(
    o_rate_area_id UUID,
    d_rate_area_id UUID,
    re_service_id UUID,
    c_id UUID,
    estimated_weight NUMERIC
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
    estimated_price NUMERIC;
BEGIN
    SELECT rip.per_unit_cents
    INTO per_unit_cents
    FROM re_intl_prices rip
    WHERE rip.origin_rate_area_id = o_rate_area_id
      AND rip.destination_rate_area_id = d_rate_area_id
      AND rip.service_id = re_service_id
      AND rip.contract_id = c_id;

    IF per_unit_cents IS NULL THEN
        RAISE EXCEPTION 'No per unit cents found for service item id: %, origin rate area: %, dest rate area: %, and contract_id: %', re_service_id, o_rate_area_id, d_rate_area_id, c_id;
    END IF;

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

    IF escalation_factor IS NULL THEN
        RAISE EXCEPTION 'Escalation factor not found for contract_id %', c_id;
    END IF;

    escalated_price := ROUND(per_unit_cents * escalation_factor / 100, 2);
    estimated_price := ROUND(escalated_price * (estimated_weight / 100), 2);

    RETURN estimated_price;
END;
$$ LANGUAGE plpgsql;

-- function that handles calculating price for IHPK & IUBPK service items, takes in:
-- origin rate area
-- re_services id
-- contract id
-- prime estimated weight
CREATE OR REPLACE FUNCTION calculate_ihpk_iubpk_price(
    o_rate_area_id UUID,
    re_service_id UUID,
    c_id UUID,
    estimated_weight NUMERIC
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
    estimated_price NUMERIC;
BEGIN
    SELECT rip.per_unit_cents
    INTO per_unit_cents
    FROM re_intl_prices rip
    WHERE rip.origin_rate_area_id = o_rate_area_id
      AND rip.service_id = re_service_id
      AND rip.contract_id = c_id;

    IF per_unit_cents IS NULL THEN
        RAISE EXCEPTION 'No per unit cents found for service item id: %, origin rate area: %, and contract_id: %', re_service_id, o_rate_area_id, c_id;
    END IF;

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

    IF escalation_factor IS NULL THEN
        RAISE EXCEPTION 'Escalation factor not found for contract_id %', contract_id;
    END IF;

    escalated_price := ROUND(per_unit_cents * escalation_factor / 100, 2);
    estimated_price := ROUND(escalated_price * (estimated_weight / 100), 2);

    RETURN estimated_price;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE PROCEDURE update_service_item_pricing(shipment_id UUID) AS
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
                IF service_code = ''POEFSC'' THEN
                    o_zip_code := get_zip_code(shipment.id, ''pickup'');
                    d_zip_code := get_zip_code(shipment.id, ''poe'');
                END IF;

                IF service_code = ''PODFSC'' THEN
                    o_zip_code := get_zip_code(shipment.id, ''destination'');
                    d_zip_code := get_zip_code(shipment.id, ''pod'');
                END IF;

                -- getting distance between the two ZIPs
                distance := get_distance(o_zip_code, d_zip_code);

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
