-- Pop prefers tables to have plural names, this table isn't used yet,
-- so we might as well rename it now before we have to worry about backwards
-- compatibility.
alter table postal_code_to_gbloc rename to postal_code_to_gblocs;

-- This view finds the GBLOC for the first shipment of each move
CREATE VIEW move_to_gbloc AS
SELECT DISTINCT ON (sh.move_id) sh.move_id AS move_id, pctg.gbloc AS gbloc
FROM mto_shipments sh
		 JOIN addresses a ON sh.pickup_address_id = a.id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code = pctg.postal_code
		 JOIN moves m ON sh.move_id = m.id
		 JOIN orders o ON m.orders_id = o.id
ORDER BY sh.move_id, sh.created_at;

-- Add id column to postal_code_to_gblocs. This is required by Pop.
alter table postal_code_to_gblocs
	add column id uuid;

-- need to circle back and update this with hardcoded IDs
update postal_code_to_gblocs
set id = uuid_generate_v4();

alter table postal_code_to_gblocs drop constraint postal_code_to_gbloc_pkey;
alter table postal_code_to_gblocs add primary key (id);
