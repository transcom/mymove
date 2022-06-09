-- We are "renaming" a couple of fields on the ppm_shipments table. This is step 6 of the process, dropping the old
-- columns

ALTER TABLE ppm_shipments
	DROP COLUMN IF EXISTS advance_requested,
	DROP COLUMN IF EXISTS advance;
