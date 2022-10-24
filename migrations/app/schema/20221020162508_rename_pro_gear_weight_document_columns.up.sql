ALTER TABLE progear_weight_tickets
	DROP empty_weight;

ALTER TABLE progear_weight_tickets
	DROP empty_document_id;

ALTER TABLE progear_weight_tickets
	DROP constructed_weight;

ALTER TABLE progear_weight_tickets
	DROP constructed_weight_document_id;

ALTER TABLE progear_weight_tickets
	RENAME full_weight TO weight;

ALTER TABLE progear_weight_tickets
	RENAME full_document_id TO document_id;

COMMENT on COLUMN progear_weight_tickets.weight IS 'Stores the weight of the the pro-gear in pounds.';
COMMENT on COLUMN progear_weight_tickets.document_id IS 'The ID of the document that is associated with the user uploads containing the pro-gear weight.';

