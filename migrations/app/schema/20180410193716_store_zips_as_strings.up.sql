ALTER TABLE tariff400ng_zip3s ALTER COLUMN zip3 TYPE varchar(3) USING to_char(zip3, 'fm000');
ALTER TABLE tariff400ng_zip5_rate_areas ALTER COLUMN zip5 TYPE varchar(5) USING to_char(zip5, 'fm00000');
