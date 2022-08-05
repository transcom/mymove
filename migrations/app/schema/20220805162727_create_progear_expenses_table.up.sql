CREATE TABLE progear_expenses (
	id uuid PRIMARY KEY,
	ppm_shipment_id uuid NOT NULL
		CONSTRAINT progear_expenses_ppm_shipment_id_fkey
	    REFERENCES ppm_shipments,
	is_own bool,
	description varchar,
	has_weight_ticket bool,
	empty_weight int,
	empty_document_id uuid NOT NULL
	    CONSTRAINT progear_expenses_empty_document_id_fkey
	    REFERENCES documents,
	full_weight int,
	full_document_id uuid NOT NULL
	    CONSTRAINT progear_expenses_full_document_id_fkey
	    REFERENCES documents,
	status ppm_document_status,
	reason varchar,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamptz
);

CREATE INDEX progear_expenses_ppm_shipment_id_idx ON progear_expenses USING hash (ppm_shipment_id);
CREATE INDEX progear_expenses_deleted_at_idx ON progear_expenses USING btree (deleted_at);

COMMENT on TABLE progear_expenses IS 'Stores pro-gear associated information and weight docs for a PPM shipment.';
COMMENT on COLUMN progear_expenses.ppm_shipment_id IS 'The ID of the PPM shipment that this pro-gear information relates to.';
COMMENT on COLUMN progear_expenses.is_own IS 'Indicates if this information is for the customer''s own progear, otherwise, it''s the spouse''s.';
COMMENT on COLUMN progear_expenses.description IS 'Stores a description of the pro-gear move.';
COMMENT on COLUMN progear_expenses.has_weight_ticket IS 'Indicates if the user has a weight ticket for their pro-gear, otherwise they have a constructed weight.';
COMMENT on COLUMN progear_expenses.empty_weight IS 'Stores the weight of the vehicle not including the pro-gear.';
COMMENT on COLUMN progear_expenses.empty_document_id IS 'The ID of the document that is associated with the user uploads containing the empty vehicle weight.';
COMMENT on COLUMN progear_expenses.full_weight IS 'Stores the weight of the vehicle including the pro-gear.';
COMMENT on COLUMN progear_expenses.full_document_id IS 'The ID of the document that is associated with the user uploads containing the full vehicle weight.';
COMMENT on COLUMN progear_expenses.status IS 'Status of the expense, e.g. APPROVED.';
COMMENT on COLUMN progear_expenses.reason IS 'Contains the reason an expense is excluded or rejected; otherwise null.';
