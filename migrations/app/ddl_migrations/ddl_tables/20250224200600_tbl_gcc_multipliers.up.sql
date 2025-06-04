-- B-23736  Daniel Jordan  adding gcc_multipliers table
CREATE TABLE IF NOT EXISTS gcc_multipliers (
    id          			uuid NOT NULL PRIMARY KEY,
	multiplier              numeric(5, 2) NOT NULL,
	start_date              date UNIQUE NOT NULL,
    end_date                date UNIQUE NOT NULL,
    created_at  			timestamp   NOT NULL DEFAULT NOW(),
    updated_at  			timestamp   NOT NULL DEFAULT NOW()
);


COMMENT ON TABLE gcc_multipliers IS 'Stores GCC multipliers to be applied to incentives';
COMMENT ON COLUMN gcc_multipliers.id IS 'Unique identifier for each multiplier record';
COMMENT ON COLUMN gcc_multipliers.multiplier IS 'The multiplier applied to the total amount, represented as a decimal (e.g., 1.30 for 30%)';
COMMENT ON COLUMN gcc_multipliers.start_date IS 'The start date for when the multiplier is active';
COMMENT ON COLUMN gcc_multipliers.end_date IS 'The end date for when the multiplier is no longer active';
COMMENT ON COLUMN gcc_multipliers.created_at IS 'Timestamp indicating when the record was created';
COMMENT ON COLUMN gcc_multipliers.updated_at IS 'Timestamp indicating when the record was last updated';