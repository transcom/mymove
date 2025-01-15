CREATE OR REPLACE FUNCTION calculate_ppm_incentive(
    ppm_id UUID,
    pickup_address_id UUID,
    destination_address_id UUID,
    move_date DATE,
    mileage INT,
    weight INT,
    is_estimated BOOLEAN,
    is_actual BOOLEAN,
    is_max BOOLEAN
) RETURNS TABLE (
    total_incentive NUMERIC,
    price_islh NUMERIC,
    price_ihpk NUMERIC,
    price_ihupk NUMERIC,
    price_fsc NUMERIC
) AS
$$
DECLARE
    ppm RECORD;
    contract_id UUID;
    o_rate_area_id UUID;
    d_rate_area_id UUID;
    service_id UUID;
    estimated_fsc_multiplier NUMERIC;
    fuel_price NUMERIC;
    price_difference NUMERIC;
    cents_above_baseline NUMERIC;
BEGIN

    IF NOT is_estimated AND NOT is_actual AND NOT is_max THEN
        RAISE EXCEPTION 'is_estimated, is_actual, and is_max cannot all be FALSE. No update will be performed.';
    END IF;

    -- Validating it's a real PPM
    SELECT ppms.id INTO ppm FROM ppm_shipments ppms WHERE ppms.id = ppm_id;
    IF ppm IS NULL THEN
        RAISE EXCEPTION 'PPM with ID % not found', ppm_id;
    END IF;

    contract_id := get_contract_id(move_date);
    IF contract_id IS NULL THEN
        RAISE EXCEPTION 'Contract not found for date: %', move_date;
    END IF;

    o_rate_area_id := get_rate_area_id(pickup_address_id, NULL, contract_id);
    IF o_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Origin rate area is NULL for address ID %', pickup_address_id;
    END IF;

    d_rate_area_id := get_rate_area_id(destination_address_id, NULL, contract_id);
    IF d_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Destination rate area is NULL for address ID %', destination_address_id;
    END IF;

    -- ISLH calculation
    SELECT rs.id INTO service_id FROM re_services rs WHERE rs.code = 'ISLH';
    price_islh := ROUND(
        calculate_escalated_price(
            o_rate_area_id,
            d_rate_area_id,
            service_id,
            contract_id,
            'ISLH',
            move_date
        ) * (weight / 100)::NUMERIC * 100, 0
    );

    -- IHPK calculation
    SELECT rs.id INTO service_id FROM re_services rs WHERE rs.code = 'IHPK';
    price_ihpk := ROUND(
        calculate_escalated_price(
            o_rate_area_id,
            NULL,
            service_id,
            contract_id,
            'IHPK',
            move_date
        ) * (weight / 100)::NUMERIC * 100, 0
    );

    -- IHUPK calculation
    SELECT rs.id INTO service_id FROM re_services rs WHERE rs.code = 'IHUPK';
    price_ihupk := ROUND(
        calculate_escalated_price(
            NULL,
            d_rate_area_id,
            service_id,
            contract_id,
            'IHUPK',
            move_date
        ) * (weight / 100)::NUMERIC * 100, 0
    );

    -- FSC calculation
    estimated_fsc_multiplier := get_fsc_multiplier(weight);
    fuel_price := get_fuel_price(move_date);
    price_difference := calculate_price_difference(fuel_price);
    cents_above_baseline := mileage * estimated_fsc_multiplier;
    price_fsc := ROUND((cents_above_baseline * price_difference) * 100);

    -- Total incentive
    total_incentive := price_islh + price_ihpk + price_ihupk + price_fsc;

    -- Update the PPM incentive values
    UPDATE ppm_shipments
    SET estimated_incentive = CASE WHEN is_estimated THEN total_incentive ELSE estimated_incentive END,
        final_incentive = CASE WHEN is_actual THEN total_incentive ELSE final_incentive END,
        max_incentive = CASE WHEN is_max THEN total_incentive ELSE max_incentive END
    WHERE id = ppm_id;

    -- Return all values
    RETURN QUERY SELECT total_incentive, price_islh, price_ihpk, price_ihupk, price_fsc;
END;
$$ LANGUAGE plpgsql;
