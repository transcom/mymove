ALTER TABLE weight_ticket_set_documents
    ALTER COLUMN move_document_id SET NOT NULL,
    ALTER COLUMN trailer_ownership_missing SET NOT NULL,
    ALTER COLUMN full_weight_ticket_missing SET NOT NULL,
    ALTER COLUMN empty_weight_ticket_missing SET NOT NULL,
    ALTER COLUMN weight_ticket_set_type SET NOT NULL;