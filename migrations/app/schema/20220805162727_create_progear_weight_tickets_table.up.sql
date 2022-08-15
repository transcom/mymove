CREATE TABLE progear_weight_tickets (
	id uuid PRIMARY KEY,
	ppm_shipment_id uuid NOT NULL
		CONSTRAINT progear_weight_tickets_ppm_shipment_id_fkey
	    REFERENCES ppm_shipments,
	belongs_to_self bool,
	description varchar,
	has_weight_tickets bool,
	empty_weight int,
	empty_document_id uuid NOT NULL
	    CONSTRAINT progear_weight_tickets_empty_document_id_fkey
	    REFERENCES documents,
	full_weight int,
	full_document_id uuid NOT NULL
	    CONSTRAINT progear_weight_tickets_full_document_id_fkey
	    REFERENCES documents,
	constructed_weight int,
	constructed_weight_document_id uuid NOT NULL
	    CONSTRAINT progear_weight_tickets_constructed_weight_document_id_fkey
	    REFERENCES documents,
	status ppm_document_status,
	reason varchar,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamptz
);

CREATE INDEX progear_weight_tickets_ppm_shipment_id_idx ON progear_weight_tickets USING hash (ppm_shipment_id);
CREATE INDEX progear_weight_tickets_deleted_at_idx ON progear_weight_tickets USING btree (deleted_at);

COMMENT on TABLE progear_weight_tickets IS 'Stores pro-gear associated information and weight docs for a PPM shipment.';
COMMENT on COLUMN progear_weight_tickets.ppm_shipment_id IS 'The ID of the PPM shipment that this pro-gear information relates to.';
COMMENT on COLUMN progear_weight_tickets.belongs_to_self IS 'Indicates if this information is for the customer''s own progear, otherwise, it''s the spouse''s.';
COMMENT on COLUMN progear_weight_tickets.description IS 'Stores a description of the pro-gear that was moved.';
COMMENT on COLUMN progear_weight_tickets.has_weight_tickets IS 'Indicates if the user has a weight ticket for their pro-gear, otherwise they have a constructed weight.';
COMMENT on COLUMN progear_weight_tickets.empty_weight IS 'Stores the weight of the vehicle not including the pro-gear.';
COMMENT on COLUMN progear_weight_tickets.empty_document_id IS 'The ID of the document that is associated with the user uploads containing the empty vehicle weight.';
COMMENT on COLUMN progear_weight_tickets.full_weight IS 'Stores the weight of the vehicle including the pro-gear.';
COMMENT on COLUMN progear_weight_tickets.full_document_id IS 'The ID of the document that is associated with the user uploads containing the full vehicle weight.';
COMMENT on COLUMN progear_weight_tickets.constructed_weight IS 'Stores the constructed weight of the pro-gear.';
COMMENT on COLUMN progear_weight_tickets.constructed_weight_document_id IS 'The ID of the document that is associated with the user uploads containing the constructed weight.';
COMMENT on COLUMN progear_weight_tickets.status IS 'Status of the expense, e.g. APPROVED.';
COMMENT on COLUMN progear_weight_tickets.reason IS 'Contains the reason an expense is excluded or rejected; otherwise null.';
