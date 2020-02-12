-- Rate areas should have the 'US' prefix.
UPDATE tariff400ng_zip5_rate_areas SET rate_area = concat('US', rate_area);
UPDATE tariff400ng_zip3s SET rate_area = concat('US', rate_area) where rate_area != 'ZIP';

-- Regions should not have the 'REGION ' prefix.
UPDATE traffic_distribution_lists SET destination_region = substring(destination_region from 8 for 20);

-- Regions should be stored as text.
ALTER table tariff400ng_zip3s ALTER region TYPE text;

-- Service areas should be stored as text.
ALTER table tariff400ng_zip3s ALTER service_area TYPE text;
ALTER table tariff400ng_service_areas ALTER service_area TYPE text;
