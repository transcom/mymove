ALTER TABLE duty_stations
	ALTER COLUMN affiliation DROP NOT NULL;

ALTER TABLE orders
	ADD COLUMN gbloc VARCHAR;
COMMENT ON COLUMN orders.gbloc IS 'TIO and TOO users from transportation offices in this GBLOC will see these orders in their queue.';
