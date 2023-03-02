-- Restore the gbloc column with the previous settings
ALTER TABLE orders
	ADD COLUMN gbloc VARCHAR;

CREATE INDEX orders_gbloc_idx ON orders (gbloc);

COMMENT ON COLUMN orders.gbloc IS 'Services Counselor and TIO users from transportation offices in this GBLOC will see these orders in their queue.';


-- Backfill gbloc data
-- Start with USMC cases
UPDATE orders
SET gbloc = 'USMC'
WHERE orders.service_member_id IN (
    SELECT id FROM service_members WHERE affiliation = 'MARINES'
);

-- Update the remaining orders records. Use the query currently used to define origin_duty_location_to_gbloc view
UPDATE orders
SET gbloc = (
    SELECT pctg.gbloc
    FROM duty_locations dl
		 JOIN addresses a ON a.id = dl.address_id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code::TEXT = pctg.postal_code::TEXT
    WHERE dl.id = orders.origin_duty_location_id
)
WHERE orders.gbloc IS NULL;

-- Drop origin_duty_location_to_gbloc view
DROP VIEW IF EXISTS origin_duty_location_to_gbloc