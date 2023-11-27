-- adding column to uploads table so we can determine if an upload is a weight ticket or not
-- this will be a boolean data type
ALTER TABLE proof_of_service_docs
ADD COLUMN is_weight_ticket boolean DEFAULT false NOT NULL;

-- Column comments
COMMENT ON COLUMN proof_of_service_docs.is_weight_ticket IS 'Determines if the proof of service doc is a weight ticket or not, this will be used in the UI when reviewing the requests.';