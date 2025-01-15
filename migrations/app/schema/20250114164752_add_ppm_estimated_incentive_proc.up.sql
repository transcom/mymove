-- function that calculates a ppm incentive given mileage, weight, and dates
-- this is used to calculate estimated, max, and actual incentives
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
) RETURNS NUMERIC AS
$$
DECLARE
    ppm RECORD;
    escalated_price NUMERIC;
    estimated_price_islh NUMERIC;
    estimated_price_ihpk NUMERIC;
    estimated_price_ihupk NUMERIC;
    estimated_price_fsc NUMERIC;
    total_incentive NUMERIC := 0;
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

    -- validating it's a real PPM
    SELECT ppms.id
    INTO ppm
    FROM ppm_shipments ppms
    WHERE ppms.id = ppm_id;

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
    estimated_price_islh := ROUND(
        calculate_escalated_price(
            o_rate_area_id,
            d_rate_area_id,
            service_id,
            contract_id,
            'ISLH',
            move_date
        ) * (weight / 100)::NUMERIC * 100, 0
    );
    RAISE NOTICE 'Estimated price for ISLH: % cents', estimated_price_islh;

    -- IHPK calculation
    SELECT rs.id INTO service_id FROM re_services rs WHERE rs.code = 'IHPK';
    estimated_price_ihpk := ROUND(
        calculate_escalated_price(
            o_rate_area_id,
            NULL,
            service_id,
            contract_id,
            'IHPK',
            move_date
        ) * (weight / 100)::NUMERIC * 100, 0
    );
    RAISE NOTICE 'Estimated price for IHPK: % cents', estimated_price_ihpk;

    -- IHUPK calculation
    SELECT rs.id INTO service_id FROM re_services rs WHERE rs.code = 'IHUPK';
    estimated_price_ihupk := ROUND(
        calculate_escalated_price(
            NULL,
            d_rate_area_id,
            service_id,
            contract_id,
            'IHUPK',
            move_date
        ) * (weight / 100)::NUMERIC * 100, 0
    );
    RAISE NOTICE 'Estimated price for IHUPK: % cents', estimated_price_ihupk;

    -- FSC calculation
    estimated_fsc_multiplier := get_fsc_multiplier(weight);
    fuel_price := get_fuel_price(move_date);
    price_difference := calculate_price_difference(fuel_price);
    cents_above_baseline := mileage * estimated_fsc_multiplier;
    estimated_price_fsc := ROUND((cents_above_baseline * price_difference) * 100);
    RAISE NOTICE 'Estimated price for FSC: % cents', estimated_price_fsc;

    -- total
    total_incentive := estimated_price_islh + estimated_price_ihpk + estimated_price_ihupk + estimated_price_fsc;
    RAISE NOTICE 'Total PPM Incentive: % cents', total_incentive;

    -- now update the incentive value
    UPDATE ppm_shipments
    SET estimated_incentive = CASE WHEN is_estimated THEN total_incentive ELSE estimated_incentive END,
        final_incentive = CASE WHEN is_actual THEN total_incentive ELSE final_incentive END,
        max_incentive = CASE WHEN is_max THEN total_incentive ELSE max_incentive END
    WHERE id = ppm_id;

    RETURN total_incentive;
END;
$$ LANGUAGE plpgsql;
