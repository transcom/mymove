-- removing valid_loa_for_tac column from transportation_accounting_codes table
ALTER TABLE transportation_accounting_codes
DROP COLUMN IF EXISTS valid_loa_for_tac;

-- then adding valid_loa_for_tac column to lines_of_accounting table
ALTER TABLE lines_of_accounting
ADD COLUMN IF NOT EXISTS valid_loa_for_tac boolean DEFAULT NULL;

-- Column comments
COMMENT ON COLUMN lines_of_accounting.valid_loa_for_tac IS 'Result of LOA service object validation that occurs whenever new LOAs are fetched by the client.';