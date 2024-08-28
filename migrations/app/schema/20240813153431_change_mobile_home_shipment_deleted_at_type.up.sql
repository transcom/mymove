ALTER TABLE mobile_homes
ALTER COLUMN deleted_at TYPE timestamp with time zone,
ALTER COLUMN deleted_at DROP NOT NULL;