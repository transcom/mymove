-- IDFSIT PerUnitCents
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
    ('fb7925e7-ebfe-49d9-9cf4-7219e68ec686'::uuid,'bd6064ca-e780-4ab4-a37b-0ae98eebb244','597bb77e-0ce7-4ba2-9624-24300962625f','2024-01-17 15:55:50.041957','2024-01-17 15:55:50.041957',false); -- PerUnitCents

-- IDASIT PerUnitCents
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
    ('51393ee1-f505-4f7b-96c4-135f771af814'::uuid,'806c6d59-57ff-4a3f-9518-ebf29ba9cb10','597bb77e-0ce7-4ba2-9624-24300962625f','2024-01-17 15:55:50.041957','2024-01-17 15:55:50.041957',false); -- PerUnitCents

-- IOFSIT PerUnitCents
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
    ('7518ec84-0c40-4c17-86dd-3ce04e2fe701'::uuid,'b488bf85-ea5e-49c8-ba5c-e2fa278ac806','597bb77e-0ce7-4ba2-9624-24300962625f','2024-01-17 15:55:50.041957','2024-01-17 15:55:50.041957',false); -- PerUnitCents

-- IOASIT PerUnitCents
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
    ('cff34123-e2a5-40ed-9cf3-451701850a26'::uuid,'bd424e45-397b-4766-9712-de4ae3a2da36','597bb77e-0ce7-4ba2-9624-24300962625f','2024-01-17 15:55:50.041957','2024-01-17 15:55:50.041957',false); -- PerUnitCents

-- inserting PortZip param for FSC
-- we need this for international PPMs since they only get reimbursed for the CONUS -> Port portion
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
	 ('bb53e034-80c2-420e-8492-f54d2018fff1'::uuid,'4780b30c-e846-437a-b39a-c499a6b09872','d9ad3878-4b94-4722-bbaf-d4b8080f339d','2024-01-17 15:55:50.041957','2024-01-17 15:55:50.041957',true); -- PortZip

-- remove PriceAreaIntlOrigin, we don't need it
DELETE FROM service_params
WHERE service_item_param_key_id = '6d44624c-b91b-4226-8fcd-98046e2f433d';

DELETE FROM service_item_param_keys
WHERE key = 'PriceAreaIntlOrigin';

-- remove PriceAreaIntlDest, we don't need it
DELETE FROM service_params
WHERE service_item_param_key_id = '4736f489-dfda-4df1-a303-8c434a120d5d';

DELETE FROM service_item_param_keys
WHERE key = 'PriceAreaIntlDest';

-- adding port info that PPMs will consume
INSERT INTO public.ports
(id, port_code, port_type, port_name, created_at, updated_at)
VALUES('d8776c6b-bc5e-45d8-ac50-ab60c34c022d'::uuid, '4E1', 'S','TACOMA, PUGET SOUND', now(), now());

INSERT INTO public.port_locations
(id, port_id, cities_id, us_post_region_cities_id, country_id, is_active, created_at, updated_at)
VALUES('ee3a97dc-112e-4805-8518-f56f2d9c6cc6'::uuid, 'd8776c6b-bc5e-45d8-ac50-ab60c34c022d'::uuid, 'baaf6ab1-6142-4fb7-b753-d0a142c75baf'::uuid, '86fef297-d61f-44ea-afec-4f679ce686b7'::uuid, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, true, now(), now());

-- func to fetch a service id from re_services by providing the service code
CREATE OR REPLACE FUNCTION get_service_id(service_code TEXT) RETURNS UUID AS $$
DECLARE
    service_id UUID;
BEGIN
    SELECT rs.id INTO service_id FROM re_services rs WHERE rs.code = service_code;
    IF service_id IS NULL THEN
        RAISE EXCEPTION 'Service code % not found in re_services', service_code;
    END IF;
    RETURN service_id;
END;
$$ LANGUAGE plpgsql;


-- db func that will calculate a PPM's incentives
-- this is used for estimated/final/max incentives
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

    -- validating it's a real PPM
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
    service_id := get_service_id('ISLH');
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
    service_id := get_service_id('IHPK');
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
    service_id := get_service_id('IHUPK');
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

    total_incentive := price_islh + price_ihpk + price_ihupk + price_fsc;

    UPDATE ppm_shipments
    SET estimated_incentive = CASE WHEN is_estimated THEN total_incentive ELSE estimated_incentive END,
        final_incentive = CASE WHEN is_actual THEN total_incentive ELSE final_incentive END,
        max_incentive = CASE WHEN is_max THEN total_incentive ELSE max_incentive END
    WHERE id = ppm_id;

    -- returning a table so we can use this data in the breakdown for the service member
    RETURN QUERY SELECT total_incentive, price_islh, price_ihpk, price_ihupk, price_fsc;
END;
$$ LANGUAGE plpgsql;


-- db func that will calculate a PPM's SIT cost
-- returns a table with total cost and the cost of each first day/add'l day SIT service item
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
            move_date
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
            move_date
        ) * (weight / 100)::NUMERIC * 100 * sit_days
    )::INT;

    -- add em up
    total_cost := price_first_day + price_addl_day;

    RETURN QUERY SELECT total_cost, price_first_day, price_addl_day;
END;
$$ LANGUAGE plpgsql;

