-- this table will be used to hold validation codes a customer will enter prior to beginning their move

CREATE TABLE IF NOT EXISTS application_parameters (
	id uuid PRIMARY KEY NOT NULL,
	validation_code TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE application_parameters IS 'Table to hold validation codes that will validate customers ability to begin their moves.';
COMMENT ON COLUMN application_parameters.validation_code IS 'Validation code, alphanumeric value that can be up to 20 characters';
