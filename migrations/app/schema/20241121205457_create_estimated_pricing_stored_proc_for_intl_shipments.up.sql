ALTER TABLE addresses
	ADD CONSTRAINT us_post_region_cities_id_fkey
	FOREIGN KEY (us_post_region_cities_id) REFERENCES us_post_region_cities (id);

DROP FUNCTION IF EXISTS calculate_escalated_price(uuid, uuid, uuid, uuid);

-- function to get the rate area id
CREATE OR REPLACE FUNCTION get_rate_area_id(
    address_id UUID,
    service_item_id UUID,
    OUT o_rate_area_id UUID
)
RETURNS UUID AS $$
BEGIN
    SELECT ro.rate_area_id
    INTO o_rate_area_id
    FROM addresses a
    JOIN re_oconus_rate_areas ro
    ON a.us_post_region_cities_id = ro.us_post_region_cities_id
    WHERE a.id = address_id;

    IF o_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate area not found for address % for service item id %', address_id, service_item_id;
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
    c_id UUID
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
BEGIN
    SELECT intl_prices.per_unit_cents
    INTO per_unit_cents
    FROM re_intl_prices intl_prices
    WHERE intl_prices.origin_rate_area_id = o_rate_area_id
      AND intl_prices.destination_rate_area_id = d_rate_area_id
      AND intl_prices.service_id = re_service_id
      AND intl_prices.contract_id = c_id;

    -- IF per_unit_cents IS NULL THEN
    --     RAISE EXCEPTION 'No matching price found for the given parameters';
    -- END IF;

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

    IF escalation_factor IS NULL THEN
        RAISE EXCEPTION 'Escalation factor not found for contract_id %', contract_id;
    END IF;

    -- Calculate the escalated price
    escalated_price := ROUND(per_unit_cents * escalation_factor::NUMERIC / 100, 2);

    RETURN escalated_price;
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

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

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

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

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
BEGIN
    SELECT ms.id, ms.pickup_address_id, ms.destination_address_id, ms.requested_pickup_date, ms.prime_estimated_weight
    INTO shipment
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF shipment IS NULL THEN
        RAISE EXCEPTION ''Shipment with ID % not found'', shipment_id;
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
                escalated_price := calculate_escalated_price(o_rate_area_id, d_rate_area_id, service_item.re_service_id, contract_id);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND(escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC, 2);
                END IF;

            WHEN service_code IN (''IHPK'', ''IUBPK'') THEN
                -- perform IHPK/IUBPK-specific logic (no destination rate area)
                o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id);
                contract_id := get_contract_id(shipment.requested_pickup_date);
                escalated_price := calculate_escalated_price(o_rate_area_id, NULL, service_item.re_service_id, contract_id);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND(escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC, 2);
                END IF;

            WHEN service_code IN (''IHUPK'', ''IUBUPK'') THEN
                -- perform IHUPK/IUBUPK-specific logic (no origin rate area)
                d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id);
                contract_id := get_contract_id(shipment.requested_pickup_date);
                escalated_price := calculate_escalated_price(NULL, d_rate_area_id, service_item.re_service_id, contract_id);

                IF shipment.prime_estimated_weight IS NOT NULL THEN
                    estimated_price := ROUND(escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC, 2);
                END IF;

            WHEN service_code IN (''POEFSC'', ''PODFSC'') THEN
                estimated_price := 0; -- placeholder
        END CASE;

        -- Update the pricing_estimate in mto_service_items
        UPDATE mto_service_items
        SET pricing_estimate = estimated_price
        WHERE id = service_item.id;
    END LOOP;
END;
'
LANGUAGE plpgsql;
