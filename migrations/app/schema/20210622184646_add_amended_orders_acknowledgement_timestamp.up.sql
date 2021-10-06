ALTER TABLE orders
    ADD COLUMN amended_orders_acknowledged_at timestamp;

COMMENT ON COLUMN orders.amended_orders_acknowledged_at IS 'A timestamp that captures when new amended orders are reviewed after a move was previously approved with original orders';
