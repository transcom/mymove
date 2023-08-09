-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

-- We had to split this migration into multiple parts because it was too large for the virus scanner.
-- Locally, we only create a few records, so we don't need to load any data in parts 2 and 3 of this migration.

-- Create indices after loading data
CREATE INDEX transportation_accounting_codes_tac_idx ON transportation_accounting_codes (tac);

ALTER TABLE transportation_accounting_codes
	ADD CONSTRAINT transportation_acounting_codes_pkey PRIMARY KEY (id);
