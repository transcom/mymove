--B-22463 M.Inthavongsay added get_application_parameter_value
CREATE OR REPLACE FUNCTION get_application_parameter_value(param_name VARCHAR)
RETURNS VARCHAR AS $$
    DECLARE param_value VARCHAR;
BEGIN

    SELECT ap.parameter_value
    INTO param_value
    FROM application_parameters ap
    WHERE LOWER(ap.parameter_name) = LOWER(param_name);

    RETURN param_value;
END;
$$ LANGUAGE plpgsql;