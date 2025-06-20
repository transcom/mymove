-- B-23819  Jim Hawks  remove unused col. set required cols to not null.

ALTER TABLE backup_contacts
  DROP COLUMN IF EXISTS name,
  ALTER COLUMN first_name SET NOT NULL,
  ALTER COLUMN last_name SET NOT NULL;
