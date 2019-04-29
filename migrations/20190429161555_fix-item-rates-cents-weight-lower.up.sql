-- These charges are subject to a min of 1,000lbs, so that rate should apply to weights < 1,000 also
UPDATE tariff400ng_item_rates
SET weight_lbs_lower = 0
WHERE weight_lbs_lower = 1000
    AND code IN ('125A', '125C', '210A', '210D', '225A', '225B')
    AND effective_date_lower = '2019-05-15';


-- These rates were assumed to be listed in cents but they were in dollars, though they were already scaled by 10
-- because they contained decimal values
UPDATE tariff400ng_item_rates
SET rate_cents = rate_cents * 10
WHERE code IN ('125C', '210D', '225B')
	AND effective_date_lower = '2019-05-15';

-- These rates were assumed to be listed in cents but they were in dollars
UPDATE tariff400ng_item_rates
SET rate_cents = rate_cents * 100
WHERE code IN ('125A', '210A', '225A')
	AND effective_date_lower = '2019-05-15';
