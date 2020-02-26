-- Deduplicate rows in tariff400ng_zip5_rate_areas table
SELECT DISTINCT ON (zip5, rate_area) *
INTO temp_zip5s
FROM tariff400ng_zip5_rate_areas;

DROP TABLE tariff400ng_zip5_rate_areas;

ALTER TABLE temp_zip5s RENAME TO tariff400ng_zip5_rate_areas;
