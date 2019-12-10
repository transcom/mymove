ALTER TABLE customers
	ADD COLUMN first_name text NOT NULL,
	ADD COLUMN last_name text NOT NULL,
	ADD COLUMN email text NOT NULL,
	ADD COLUMN phone text NOT NULL,
	ADD COLUMN dod_id text;