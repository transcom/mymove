-- These charges are subject to a min of 1,000lbs, so that rate should apply to weights < 1,000 also
UPDATE tariff400ng_item_rates
SET weight_lbs_lower = 0
WHERE weight_lbs_lower = 1000
    AND code IN ('225A', '225B');

-- Misc charge, which is charged at a rate of $1 for every quantity (dollars)
INSERT INTO tariff400ng_item_rates (id, created_at, updated_at, code, rate_cents, effective_date_lower, effective_date_upper) VALUES ('50d6bfbb-77a8-4697-b1b1-68005d0f6569', now(), now(), '226A', 100, '2018-05-15', '2019-05-15');
