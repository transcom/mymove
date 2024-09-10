ALTER TYPE upload_type ADD VALUE 'APP';
COMMENT ON COLUMN uploads.upload_type IS 'Who created the upload: USER, PRIME, OFFICE, or APPLICATION';