-- This constraint is a remnant of the old TGET system and is no longer needed.
-- The old TGET system was based off hard coded imports and FK setting, this `loa_id`
-- column is NOT TGET data, we don't need this. We lookup based on composite key.
ALTER TABLE transportation_accounting_codes
DROP CONSTRAINT transportation_accounting_codes_loa_id_fkey;