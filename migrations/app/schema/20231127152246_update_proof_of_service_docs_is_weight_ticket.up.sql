-- Since we've added is_weight_ticket as a default FALSE this will go ahead and set all existing NULLs to false
UPDATE proof_of_service_docs
SET is_weight_ticket = FALSE
WHERE is_weight_ticket IS NULL;
