ALTER TABLE tariff400ng_zip3s ALTER COLUMN zip3 TYPE integer USING (zip3::integer);
ALTER TABLE tariff400ng_zip5_rate_areas ALTER COLUMN zip5 TYPE integer USING (zip5::integer);
