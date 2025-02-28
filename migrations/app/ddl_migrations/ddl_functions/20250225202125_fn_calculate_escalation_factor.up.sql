-- Fetch the escalation factor given a contract ID and date
CREATE OR REPLACE FUNCTION calculate_escalation_factor(in_contract_id UUID, in_date DATE) RETURNS NUMERIC(6, 5) AS $$
DECLARE declared_factor NUMERIC(6, 5);
BEGIN
SELECT rcy.escalation INTO declared_factor
FROM re_contract_years rcy
WHERE rcy.contract_id = in_contract_id
    AND in_date >= rcy.start_date
    AND in_date <= rcy.end_date
ORDER BY rcy.start_date DESC
LIMIT 1;
IF NOT FOUND THEN RAISE EXCEPTION 'No matching contract year found for contract_id=% and date=%',
in_contract_id,
in_date;
END IF;
RETURN declared_factor;
END;
$$ LANGUAGE plpgsql;