ALTER TABLE application_parameters
ADD COLUMN IF NOT EXISTS parameter_name TEXT,
ADD COLUMN IF NOT EXISTS parameter_value TEXT,
ALTER COLUMN validation_code DROP NOT NULL;

COMMENT ON COLUMN application_parameters.parameter_name IS 'The name of the parameter';
COMMENT ON COLUMN application_parameters.parameter_value IS 'The value of the parameter';
