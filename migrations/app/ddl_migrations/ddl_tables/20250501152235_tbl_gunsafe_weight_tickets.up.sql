--B-23342 Tae Jung create gunsafe_weight_tickets table for gun safe feature E-06078
CREATE TABLE IF NOT EXISTS gunsafe_weight_tickets (
	id uuid PRIMARY KEY,
	ppm_shipment_id uuid NOT NULL
		CONSTRAINT gunsafe_weight_tickets_ppm_shipment_id_fkey
	    REFERENCES ppm_shipments,
	has_weight_tickets bool,
	description varchar,
	weight int CHECK (weight >= 0),
	document_id uuid NOT NULL
	    CONSTRAINT gunsafe_weight_tickets_document_id_fkey
	    REFERENCES documents,
	status ppm_document_status,
	reason varchar,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamptz
);

-- B-23368 need these "submitted" columns to stay consistent with other PPM documents
ALTER TABLE gunsafe_weight_tickets
ADD COLUMN IF NOT EXISTS submitted_has_weight_tickets bool DEFAULT false;

ALTER TABLE gunsafe_weight_tickets
ADD COLUMN IF NOT EXISTS submitted_weight int CHECK (weight >= 0);

CREATE INDEX IF NOT EXISTS gunsafe_weight_tickets_ppm_shipment_id_idx ON gunsafe_weight_tickets USING hash (ppm_shipment_id);
CREATE INDEX IF NOT EXISTS gunsafe_weight_tickets_deleted_at_idx ON gunsafe_weight_tickets USING btree (deleted_at);

COMMENT on TABLE gunsafe_weight_tickets IS 'Stores gun safe associated information and weight docs for a PPM shipment.';
COMMENT on COLUMN gunsafe_weight_tickets.ppm_shipment_id IS 'The ID of the PPM shipment that this gun safe information relates to.';
COMMENT on COLUMN gunsafe_weight_tickets.has_weight_tickets IS 'Indicates if the user has a weight ticket for their gun safe.';
COMMENT on COLUMN gunsafe_weight_tickets.description IS 'Stores a description of the gun safe that was moved.';
COMMENT on COLUMN gunsafe_weight_tickets.weight IS 'Stores the weight of the gun safe.';
COMMENT on COLUMN gunsafe_weight_tickets.document_id IS 'The ID of the document that is associated with the user uploads containing the gun safe weight.';
COMMENT on COLUMN gunsafe_weight_tickets.status IS 'Status of the expense, e.g. APPROVED.';
COMMENT on COLUMN gunsafe_weight_tickets.reason IS 'Contains the reason an expense is excluded or rejected; otherwise null.';
COMMENT on COLUMN gunsafe_weight_tickets.submitted_weight IS 'Stores the customer submitted value for the weight of the gun safe when uploading weight tickets';
COMMENT on COLUMN gunsafe_weight_tickets.submitted_has_weight_tickets IS 'Stores the customer submitted value for the "I don''t have this weight ticket" checkbox when uploading gun safe weights';