-- This migration removes unused payment request status type of RECEIVED_BY_GEX
-- all previous payment requests using type were updated to TPPS_RECEIVED in
-- migrations/app/schema/20240725190050_update_payment_request_status_tpps_received.up.sql

-- update again in case new payment requests have used this status
UPDATE payment_requests SET status = 'TPPS_RECEIVED' where status = 'RECEIVED_BY_GEX';

--- rename existing enum
ALTER TYPE payment_request_status RENAME TO payment_request_status_temp;

-- create a new enum with both old and new statuses - both old and new statuses must exist in the enum to do the update setting old to new
CREATE TYPE payment_request_status AS ENUM(
    'PENDING',
    'REVIEWED',
    'SENT_TO_GEX',
    'PAID',
    'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED',
    'EDI_ERROR',
    'DEPRECATED',
    'TPPS_RECEIVED'
    );

alter  table payment_requests alter column  status drop default;
alter  table payment_requests alter column  status drop not null;

-- alter the payment_requests status column to use the new enum
ALTER TABLE payment_requests ALTER COLUMN status TYPE payment_request_status USING status::text::payment_request_status;

-- get rid of the temp type
DROP TYPE payment_request_status_temp;

ALTER TABLE payment_requests
ALTER COLUMN status SET DEFAULT 'PENDING',
ALTER COLUMN status SET NOT NULL;