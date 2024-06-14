ALTER TABLE weight_tickets
    ADD COLUMN IF NOT EXISTS submitted_owns_trailer boolean DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS submitted_trailer_meets_criteria boolean DEFAULT NULL;

COMMENT ON COLUMN weight_tickets.submitted_owns_trailer IS 'Stores the customer submitted owns_trailer.';
COMMENT ON COLUMN weight_tickets.submitted_trailer_meets_criteria IS 'Stores the customer submitted trailer_meets_criteria.';

ALTER TABLE progear_weight_tickets
    ADD COLUMN IF NOT EXISTS submitted_belongs_to_self boolean DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS submitted_has_weight_tickets boolean DEFAULT NULL;

COMMENT ON COLUMN progear_weight_tickets.submitted_belongs_to_self IS 'Stores the customer belongs_to_self.';
COMMENT ON COLUMN progear_weight_tickets.submitted_has_weight_tickets IS 'Stores the customer submitted has_weight_tickets.';
