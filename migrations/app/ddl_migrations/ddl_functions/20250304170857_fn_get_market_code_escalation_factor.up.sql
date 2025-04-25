-- B-22662 C.Jewell initial shipment type table market factor escalation lookup
CREATE OR REPLACE FUNCTION get_market_code_factor_escalation(
        in_rate_area_id UUID,
        in_contract_id UUID,
        in_re_service_id UUID
    ) RETURNS NUMERIC LANGUAGE plpgsql AS $$
DECLARE declared_is_oconus BOOLEAN;
declared_market_code TEXT;
declared_service_code TEXT;
declared_factor NUMERIC;
BEGIN
-- Currently only INPK is configured for market factors
SELECT rs.code INTO declared_service_code
FROM re_services rs
WHERE rs.id = in_re_service_id
LIMIT 1;
-- Catch unsupported codes. In the future it is likely we may support more
IF declared_service_code IS NULL THEN RAISE EXCEPTION 'No re_services found for id=% while fetching market factor escalation',
in_re_service_id;
ELSIF declared_service_code != 'INPK' THEN RAISE EXCEPTION 'get_market_code_factor is only for INPK, but got code=% in_re_service_id=%',
declared_service_code,
in_re_service_id;
END IF;
-- Check the market of the rate area
SELECT ra.is_oconus INTO declared_is_oconus
FROM re_rate_areas ra
WHERE ra.id = in_rate_area_id
LIMIT 1;
IF declared_is_oconus THEN declared_market_code := 'O';
ELSE declared_market_code := 'C';
END IF;
-- Fetch the market factor
SELECT stp.factor INTO declared_factor
FROM re_shipment_type_prices stp
WHERE stp.contract_id = in_contract_id
    AND stp.service_id = in_re_service_id
    AND stp.market = declared_market_code
LIMIT 1;
IF declared_factor IS NULL THEN RAISE EXCEPTION 'No OCONUS/CONUS market factor found market=% in_contract_id=% service_id=%',
declared_market_code,
in_contract_id,
in_re_service_id;
END IF;
RETURN declared_factor;
END;
$$;