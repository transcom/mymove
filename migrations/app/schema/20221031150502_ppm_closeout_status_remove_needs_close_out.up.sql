--- rename existing enum
ALTER TYPE ppm_shipment_status RENAME TO ppm_shipment_status_temp;

-- create a new enum without the value we're looking to remove (NEEDS_CLOSE_OUT)
CREATE TYPE ppm_shipment_status AS ENUM('DRAFT', 'SUBMITTED', 'WAITING_ON_CUSTOMER', 'NEEDS_ADVANCE_APPROVAL', 'NEEDS_PAYMENT_APPROVAL', 'PAYMENT_APPROVED');

-- Remove references to the old value from the ppm_shipments table (there probably aren't any, but this is for safety)
UPDATE ppm_shipments set status = 'NEEDS_PAYMENT_APPROVAL' WHERE status = 'NEEDS_CLOSE_OUT';

-- alter the ppm shipments status column to use the new enum
ALTER TABLE ppm_shipments ALTER COLUMN status TYPE ppm_shipment_status USING status::text::ppm_shipment_status;

-- get rid of the temp type
DROP TYPE ppm_shipment_status_temp

