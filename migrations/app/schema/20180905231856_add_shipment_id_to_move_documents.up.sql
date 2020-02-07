-- Add Shipment ID field to move_documents
ALTER TABLE move_documents
	ADD COLUMN shipment_id uuid;

-- Add foreign key constraint
ALTER TABLE move_documents
	ADD CONSTRAINT move_documents_shipment_id_fkey
	FOREIGN KEY (shipment_id) REFERENCES shipments (id)
	ON DELETE RESTRICT;
