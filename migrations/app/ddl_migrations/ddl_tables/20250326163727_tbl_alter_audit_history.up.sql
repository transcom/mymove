CREATE SEQUENCE IF NOT EXISTS audit_seq START WITH 1;

ALTER TABLE audit_history ADD COLUMN IF NOT EXISTS seq_num serial;

ALTER TABLE audit_history ALTER COLUMN seq_num SET DEFAULT nextval('audit_seq');
