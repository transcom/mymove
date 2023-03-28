-- Restore the gbloc column with the previous settings
ALTER TABLE orders
	ADD COLUMN gbloc VARCHAR;

COMMENT ON COLUMN orders.gbloc IS 'Services Counselor office users from transportation offices in this GBLOC will see these orders in their queue.';


-- Backfill gbloc data
UPDATE orders
SET gbloc = (
    SELECT pctg.gbloc
    FROM duty_locations dl
		 JOIN addresses a ON a.id = dl.address_id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code::TEXT = pctg.postal_code::TEXT
    WHERE dl.id = orders.origin_duty_location_id
)
WHERE orders.gbloc IS NULL;

-- Apply index AFTER data backfill is complete
CREATE INDEX orders_gbloc_idx ON orders (gbloc);

-- Drop origin_duty_location_to_gbloc view
DROP VIEW IF EXISTS origin_duty_location_to_gbloc
