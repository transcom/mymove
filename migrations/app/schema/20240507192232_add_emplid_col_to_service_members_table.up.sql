ALTER TABLE service_members ADD COLUMN IF NOT EXISTS emplid TEXT UNIQUE DEFAULT NULL;

COMMENT ON COLUMN service_members.emplid IS 'A Coast Guard customer''s Employee Identification Number.';