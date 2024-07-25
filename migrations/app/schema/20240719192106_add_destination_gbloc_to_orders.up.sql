ALTER TABLE orders  ADD COLUMN IF NOT EXISTS destination_gbloc VARCHAR;

CREATE INDEX orders_destination_gbloc_idx ON orders (destination_gbloc);

COMMENT ON COLUMN orders.gbloc IS 'GBLOC for move destination, will be sent to Prime so that they have accurate contact info for destination.';