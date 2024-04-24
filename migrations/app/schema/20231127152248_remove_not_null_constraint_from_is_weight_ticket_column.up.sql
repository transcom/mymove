-- dropping not null constraint to allow for it to be an optional value
ALTER TABLE proof_of_service_docs ALTER COLUMN is_weight_ticket DROP NOT NULL;