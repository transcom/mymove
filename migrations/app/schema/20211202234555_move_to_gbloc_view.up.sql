-- Pop prefers tables to have plural names, this table isn't used yet,
-- so we might as well rename it now before we have to worry about backwards
-- compatibility.
ALTER TABLE postal_code_to_gbloc
	RENAME TO postal_code_to_gblocs;

-- This view finds the GBLOC for the first shipment of each move
CREATE VIEW move_to_gbloc AS
SELECT DISTINCT ON (sh.move_id) sh.move_id AS move_id, pctg.gbloc AS gbloc
FROM mto_shipments sh
		 JOIN addresses a ON sh.pickup_address_id = a.id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code = pctg.postal_code
ORDER BY sh.move_id, sh.created_at;

-- Add id column to postal_code_to_gblocs. This is required by Pop.
ALTER TABLE postal_code_to_gblocs
	ADD COLUMN id uuid;

-- need to circle back and update this with hardcoded IDs
UPDATE postal_code_to_gblocs
SET id = uuid_generate_v4();

-- Now that we've got our new ID field populated, let's add back all the indices and constraints
ALTER TABLE postal_code_to_gblocs
	DROP CONSTRAINT postal_code_to_gbloc_pkey,
	ADD PRIMARY KEY (id),
	ALTER COLUMN postal_code SET NOT NULL,
	ADD CONSTRAINT unique_postal_code UNIQUE (postal_code);

-- We were originally planning to use this field, but decided to write a query to look the
-- GBLOC up on the fly instead. No code references this column.
ALTER TABLE orders
	DROP COLUMN gbloc;
