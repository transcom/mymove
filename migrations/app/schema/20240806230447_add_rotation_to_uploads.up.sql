ALTER TABLE uploads ADD COLUMN IF NOT EXISTS rotation INTEGER;
COMMENT ON COLUMN uploads.rotation IS 'Adjusted rotation of the doc image';
