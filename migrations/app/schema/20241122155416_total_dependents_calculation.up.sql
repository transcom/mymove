-- Set temp timeout due to large file modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;
-- Zero downtime not necessary, this is not used at this time
ALTER TABLE entitlements
DROP COLUMN IF EXISTS total_dependents; -- This column has never been used, this has been confirmed prior to the migration
-- The calculation should only ever work if dependents 12 and under or 12 and over are present
-- These fields are only present on OCONUS
-- Since we don't know the number of dependents under 12 or 12 and over on CONUS moves, we don't want to default to 0 total dependents. That'd be confusing
ALTER TABLE entitlements
ADD COLUMN total_dependents integer GENERATED ALWAYS AS (
    CASE
        WHEN dependents_under_twelve IS NULL AND dependents_twelve_and_over IS NULL THEN NULL
        ELSE COALESCE(dependents_under_twelve, 0) + COALESCE(dependents_twelve_and_over, 0)
    END
) STORED;
