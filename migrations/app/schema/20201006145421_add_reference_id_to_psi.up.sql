-- Get rid of any existing payment requests and the associated tree of data attached to them
-- (which should only exist in experimental/staging).  This allows us to start over and have
-- payment service items populated with reference IDs rather than "repairing" the ones already
-- there via a migration.
TRUNCATE TABLE payment_requests CASCADE;

-- Add the new reference ID and make sure it has a unique index.
ALTER TABLE payment_service_items
    ADD COLUMN reference_id VARCHAR(255) NOT NULL UNIQUE;

COMMENT ON COLUMN payment_service_items.reference_id IS 'Shorter ID (used by EDI) to uniquely identify this payment service item. Format is the MTO reference ID, followed by a dash, followed by enough of the payment service item ID (without dashes) to make it unique.';
