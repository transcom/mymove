-- B-23301  Jim Hawks  added columns for first name and last name

ALTER TABLE backup_contacts
  ADD COLUMN IF NOT EXISTS first_name TEXT,
  ADD COLUMN IF NOT EXISTS last_name TEXT;

COMMENT ON COLUMN backup_contacts.first_name IS 'First name of the backup contact';
COMMENT ON COLUMN backup_contacts.last_name IS 'Last name of the backup contact';

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'backup_contacts'
        AND column_name = 'name'
        AND is_nullable = 'NO'
    ) THEN
        EXECUTE 'ALTER TABLE backup_contacts ALTER COLUMN name DROP NOT NULL';
    END IF;
END $$ LANGUAGE plpgsql;
