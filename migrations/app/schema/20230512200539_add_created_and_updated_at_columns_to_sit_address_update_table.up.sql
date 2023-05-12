--Add new columns and remove reason column
ALTER TABLE sit_address_updates
ADD COLUMN created_at timestamp not null,
ADD COLUMN updated_at timestamp not null,
DROP COLUMN reason;
