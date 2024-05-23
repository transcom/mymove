ALTER TABLE application_parameters
ALTER COLUMN parameter_name DROP NOT NULL,
ALTER COLUMN parameter_value DROP NOT NULL;