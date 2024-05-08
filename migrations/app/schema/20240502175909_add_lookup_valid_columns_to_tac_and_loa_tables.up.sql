-- add column to transportation_accounting_codes table to store whether tac has a valid loa so advana is aware of any issues
-- this will be a boolean data type, null until a lookup is done and then either true for valid loa or false for invalid loa
ALTER TABLE transportation_accounting_codes
ADD COLUMN IF NOT EXISTS valid_loa_for_tac boolean DEFAULT NULL;

-- Column comments
COMMENT ON COLUMN transportation_accounting_codes.valid_loa_for_tac IS 'Result of LOA lookup for a TAC.';

-- add column to lines_of_accounting table to store whether loa has a valid hhg program code so advana is aware of any issues
-- this will be a boolean data type, null until a lookup is done and then either true for valid hhg prog code or false for invalid hhg prog code
ALTER TABLE lines_of_accounting
ADD COLUMN IF NOT EXISTS valid_hhg_code_for_loa boolean DEFAULT NULL;

-- Column comments
COMMENT ON COLUMN lines_of_accounting.valid_hhg_code_for_loa IS 'Result of HHG program code lookup for a LOA.';