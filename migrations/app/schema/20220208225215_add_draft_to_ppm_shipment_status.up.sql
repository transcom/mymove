-- Create new enum type draft for ppm_shipment_status
CREATE TYPE ppm_shipment_status_2 AS enum (
	'SUBMITTED',
	'WAITING_ON_CUSTOMER',
	'NEEDS_ADVANCE_APPROVAL',
	'NEEDS_PAYMENT_APPROVAL',
	'PAYMENT_APPROVED',
    'DRAFT'
	);

--Alter the table to use our new type
ALTER TABLE ppm_shipments
	ALTER COLUMN status TYPE ppm_shipment_status_2
		USING (status::text::ppm_shipment_status_2);
--Drop the old type
DROP TYPE ppm_shipment_status;

--Rename the type so it matches the naming of the old one
ALTER TYPE ppm_shipment_status_2 RENAME to ppm_shipment_status;
