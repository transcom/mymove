-- We are "renaming" a couple of fields on the ppm_shipments table. This is step 3 of the process, back-filling the data
-- from the old columns to the new columns

UPDATE ppm_shipments
	SET
		has_requested_advance = advance_requested,
		advance_amount_requested = advance
	WHERE advance_requested is not null;
