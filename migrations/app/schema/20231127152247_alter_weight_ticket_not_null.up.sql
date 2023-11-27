-- Make is weight ticket not null now because it used to be a required field. The next migration will then delete this attribute
ALTER TABLE proof_of_service_docs
ALTER COLUMN is_weight_ticket SET NOT NULL;
