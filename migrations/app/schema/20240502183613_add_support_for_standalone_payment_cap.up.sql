ALTER TABLE application_parameters
ADD COLUMN IF NOT EXISTS parameter_name TEXT NOT NULL;

ALTER TABLE application_parameters
RENAME COLUMN validation_code TO parameter_value;

COMMENT ON COLUMN application_parameters.parameter_name IS 'The name of the parameter';
COMMENT ON COLUMN application_parameters.parameter_value IS 'The value of the parameter';