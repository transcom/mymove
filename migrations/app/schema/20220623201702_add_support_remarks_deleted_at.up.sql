ALTER TABLE customer_support_remarks
	ADD COLUMN deleted_at timestamp WITH TIME ZONE;

CREATE INDEX customer_support_remarks_deleted_at_idx
    ON customer_support_remarks (deleted_at);

COMMENT ON COLUMN customer_support_remarks.deleted_at IS 'Date & time that the customer support remark was deleted';
