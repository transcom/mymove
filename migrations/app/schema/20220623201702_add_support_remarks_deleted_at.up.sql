ALTER TABLE customer_support_remarks
	ADD COLUMN deleted_at timestamp WITH TIME ZONE;

CREATE INDEX customer_support_remarks_deleted_at_idx
    ON customer_support_remarks (deleted_at);
