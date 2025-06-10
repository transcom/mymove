ALTER TABLE re_intl_other_prices ADD COLUMN IF NOT EXISTS is_less_50_miles bool;

COMMENT ON COLUMN re_intl_other_prices.is_less_50_miles IS 'Is less than 50 miles price';

ALTER TABLE re_intl_other_prices DROP CONSTRAINT IF EXISTS re_intl_other_prices_unique_key;

ALTER TABLE re_intl_other_prices ADD CONSTRAINT re_intl_other_prices_unique_key unique (contract_id, service_id, is_peak_period, rate_area_id, is_less_50_miles);