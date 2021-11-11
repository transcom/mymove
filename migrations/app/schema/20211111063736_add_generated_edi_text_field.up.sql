-- add two new columns to entitlements table
ALTER TABLE payment_requests ADD COLUMN generated_edi_858_text text;
COMMENT ON COLUMN payment_requests.generated_edi_858_text IS 'Once a payment request is processed and an EDI 858 is generated it is stored in this column for reference in diagnosing issues'
