--B-22742  C. Kleinjan  Updating proc call for calculate_escalated_price to use more parameters
CREATE OR REPLACE FUNCTION calculate_ppm_sit_cost(
    ppm_id UUID,
    address_id UUID,
    is_origin BOOLEAN,
    move_date DATE,
    weight INT,
    sit_days INT
) RETURNS TABLE (
    total_cost INT,
    price_first_day INT,
    price_addl_day INT
) AS
$$
DECLARE
    ppm RECORD;
    contract_id UUID;
    sit_rate_area_id UUID;
    service_id UUID;
BEGIN
    -- make sure we validate parameters
    IF sit_days IS NULL OR sit_days < 0 THEN
        RAISE EXCEPTION 'SIT days must be a positive integer. Provided value: %', sit_days;
    END IF;

    SELECT ppms.id INTO ppm FROM ppm_shipments ppms WHERE ppms.id = ppm_id;
    IF ppm IS NULL THEN
        RAISE EXCEPTION 'PPM with ID % not found', ppm_id;
    END IF;

    contract_id := get_contract_id(move_date);
    IF contract_id IS NULL THEN
        RAISE EXCEPTION 'Contract not found for date: %', move_date;
    END IF;

    sit_rate_area_id := get_rate_area_id(address_id, NULL, contract_id);
    IF sit_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate area is NULL for address ID % and contract ID %', address_id, contract_id;
    END IF;

    -- calculate first day SIT cost
    service_id := get_service_id(CASE WHEN is_origin THEN 'IOFSIT' ELSE 'IDFSIT' END);
    price_first_day := (
        calculate_escalated_price(
            CASE WHEN is_origin THEN sit_rate_area_id ELSE NULL END,
            CASE WHEN NOT is_origin THEN sit_rate_area_id ELSE NULL END,
            service_id,
            contract_id,
            CASE WHEN is_origin THEN 'IOFSIT' ELSE 'IDFSIT' END,
            move_date,
            NULL
        ) * (weight / 100)::NUMERIC * 100
    )::INT;

    -- calculate additional day SIT cost
    service_id := get_service_id(CASE WHEN is_origin THEN 'IOASIT' ELSE 'IDASIT' END);
    price_addl_day := (
        calculate_escalated_price(
            CASE WHEN is_origin THEN sit_rate_area_id ELSE NULL END,
            CASE WHEN NOT is_origin THEN sit_rate_area_id ELSE NULL END,
            service_id,
            contract_id,
            CASE WHEN is_origin THEN 'IOASIT' ELSE 'IDASIT' END,
            move_date,
            NULL
        ) * (weight / 100)::NUMERIC * 100 * sit_days
    )::INT;

    -- add em up
    total_cost := price_first_day + price_addl_day;

    RETURN QUERY SELECT total_cost, price_first_day, price_addl_day;
END;
$$ LANGUAGE plpgsql;