-- Add PPM ID field to move_documents
ALTER TABLE move_documents
	ADD COLUMN personally_procured_move_id uuid;

-- Add foreign key constraint
ALTER TABLE move_documents
	ADD CONSTRAINT move_documents_personally_procured_move_id_fkey
	FOREIGN KEY (personally_procured_move_id) REFERENCES personally_procured_moves (id)
	ON DELETE RESTRICT;

-- Loop through and backfill PPM IDs using first available PPM
DO $$
	DECLARE
		md record;
		ppmId uuid;
BEGIN
    FOR md IN SELECT * FROM move_documents
    LOOP
        SELECT ppm.id INTO ppmId from personally_procured_moves ppm WHERE ppm.move_id = md.move_id ORDER BY ppm.created_at DESC LIMIT 1;
        UPDATE move_documents SET personally_procured_move_id = ppmId WHERE id = md.id;
    END LOOP;
END$$;
