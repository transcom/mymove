--Add new columns
ALTER TABLE sit_address_updates
ADD COLUMN created_at timestamp not null,
ADD COLUMN updated_at timestamp not null;
