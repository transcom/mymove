--B-22466  M.Inthavongsay Adding initial migration file for calculate_escalated_price stored procedure using new migration process.
--B-22662  C.Jewell Replaced escalation factor select with reusable func.
--Also updating to allow IOPSIT and IDDSIT SIT service items.
-- B-22742  C. Kleinjan  Migrate function to DDL Migrations and adding the ability to get escalated price for ICRT and IUCRT

-- function to calculate the escalated price, takes in:
-- origin rate area
-- dest rate area
-- re_services id
-- contract id
-- requested pickup date
-- sit delivery miles
-- adding the is_peak_period check to refine the price query further
CREATE OR REPLACE FUNCTION calculate_escalated_price(
    o_rate_area_id UUID,
    d_rate_area_id UUID,
    re_service_id UUID,
    c_id UUID,
    service_code TEXT,
    requested_pickup_date DATE,
    sit_delivery_miles NUMERIC
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
    is_oconus BOOLEAN;
    peak_period BOOLEAN;
BEGIN
    -- we need to query the appropriate table based on the service code
    -- need to establish if the shipment is being moved during peak period
    peak_period := is_peak_period(requested_pickup_date);
    IF service_code IN ('IOSHUT','IDSHUT') THEN
		IF service_code = 'IOSHUT' THEN
        	SELECT ra.is_oconus
        	INTO is_oconus
        	FROM re_rate_areas ra
        	WHERE ra.id = o_rate_area_id;
		ELSE
			SELECT ra.is_oconus
        	INTO is_oconus
        	FROM re_rate_areas ra
        	WHERE ra.id = d_rate_area_id;
		END IF;

        SELECT rip.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_accessorial_prices rip
        WHERE
            rip.market = (CASE
                WHEN is_oconus THEN 'O'
                ELSE 'C'
			END)
          AND rip.service_id = re_service_id
          AND rip.contract_id = c_id;
    ELSIF service_code IN ('IUCRT', 'ICRT') THEN
        IF service_code = 'ICRT' THEN
        	SELECT ra.is_oconus
        	INTO is_oconus
        	FROM re_rate_areas ra
        	WHERE ra.id = o_rate_area_id;
		ELSE
			SELECT ra.is_oconus
        	INTO is_oconus
        	FROM re_rate_areas ra
        	WHERE ra.id = d_rate_area_id;
		END IF;

        SELECT rip.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_accessorial_prices rip
        WHERE
            rip.market = (CASE
                WHEN is_oconus THEN 'O'
                ELSE 'C'
			END)
          AND rip.service_id = re_service_id
          AND rip.contract_id = c_id;
    ELSIF service_code IN ('ISLH', 'UBP') THEN
        SELECT rip.per_unit_cents
        INTO per_unit_cents
        FROM re_intl_prices rip
        WHERE rip.origin_rate_area_id = o_rate_area_id AND rip.destination_rate_area_id = d_rate_area_id
          AND rip.service_id = re_service_id
          AND rip.contract_id = c_id
          AND rip.is_peak_period = peak_period;
    ELSIF service_code IN ('IOPSIT', 'IDDSIT') THEN
		IF service_code = 'IOPSIT' THEN
            SELECT riop.per_unit_cents
            INTO per_unit_cents
            FROM re_intl_other_prices riop
            WHERE riop.rate_area_id = o_rate_area_id
            AND riop.service_id = re_service_id
            AND riop.contract_id = c_id
            AND riop.is_peak_period = peak_period
            AND riop.is_less_50_miles = (sit_delivery_miles <= 50);
		ELSE
            SELECT riop.per_unit_cents
            INTO per_unit_cents
            FROM re_intl_other_prices riop
            WHERE riop.rate_area_id = d_rate_area_id
            AND riop.service_id = re_service_id
            AND riop.contract_id = c_id
            AND riop.is_peak_period = peak_period
            AND riop.is_less_50_miles = (sit_delivery_miles <= 50);
		END IF;
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

    escalation_factor := calculate_escalation_factor(
        c_id,
        requested_pickup_date
    );

    -- calculate the escalated price, return in dollars (dividing by 100)
    per_unit_cents := per_unit_cents / 100; -- putting in dollars
    escalated_price := ROUND(per_unit_cents * escalation_factor, 2); -- rounding to two decimals (100.00)

    RETURN escalated_price;
END;
$$ LANGUAGE plpgsql;