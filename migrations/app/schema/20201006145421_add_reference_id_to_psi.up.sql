-- Get rid of any existing payment requests and the associated tree of data attached to them
-- (which should only exist in experimental/staging).  This allows us to start over and have
-- payment service items populated with reference IDs rather than "repairing" the ones already
-- there via a migration.
TRUNCATE TABLE payment_requests CASCADE;

-- Add the new reference ID and make sure it has a unique index.
ALTER TABLE payment_service_items
    ADD COLUMN reference_id VARCHAR(255) NOT NULL UNIQUE;
