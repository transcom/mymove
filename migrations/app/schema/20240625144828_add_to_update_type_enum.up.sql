ALTER TYPE upload_type ADD VALUE 'OFFICE';
COMMENT ON COLUMN uploads.upload_type IS 'Who created the upload: USER, PRIME, OFFICE';