--- rename existing enum
ALTER TYPE ppm_shipment_status RENAME TO ppm_shipment_status_temp;

-- create a new enum with both old and new statuses
-- why? because both old and new statuses must exist in the enum to do the update setting old to new
CREATE TYPE ppm_shipment_status AS ENUM('DRAFT', 'SUBMITTED', 'WAITING_ON_CUSTOMER', 'NEEDS_ADVANCE_APPROVAL', 'NEEDS_CLOSEOUT', 'CLOSEOUT_COMPLETE');

-- alter the ppm shipments status column to use the new enum
ALTER TABLE ppm_shipments ALTER COLUMN status TYPE ppm_shipment_status USING status::text::ppm_shipment_status;

-- get rid of the temp type
DROP TYPE ppm_shipment_status_temp;
