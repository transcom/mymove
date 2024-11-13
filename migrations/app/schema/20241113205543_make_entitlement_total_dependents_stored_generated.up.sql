UPDATE entitlements
SET total_dependents = entitlements.dependents_under_twelve + entitlements.dependents_twelve_and_over
WHERE total_dependents IS DISTINCT FROM (entitlements.dependents_under_twelve + entitlements.dependents_twelve_and_over);

ALTER TABLE entitlements DROP COLUMN total_dependents;

ALTER TABLE entitlements
ADD COLUMN total_dependents INTEGER
GENERATED ALWAYS AS (COALESCE(entitlements.dependents_under_twelve, 0) + COALESCE(entitlements.dependents_twelve_and_over, 0)) STORED;
