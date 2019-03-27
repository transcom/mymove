ALTER TABLE signed_certifications
	ADD COLUMN personally_procured_move_id uuid REFERENCES personally_procured_moves(id),
	ADD COLUMN shipment_id uuid REFERENCES shipments(id),
	ADD COLUMN certification_type text;

