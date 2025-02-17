CREATE OR REPLACE FUNCTION square_number(input_num numeric)
RETURNS numeric AS $$
BEGIN
    RETURN input_num * input_num +15 *75;
END;
$$ LANGUAGE plpgsql;
