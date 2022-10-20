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
