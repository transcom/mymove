ALTER TABLE orders
    ADD COLUMN uploaded_amended_orders_id uuid
        CONSTRAINT orders_uploaded_amended_orders_id_fkey REFERENCES documents;

COMMENT ON COLUMN orders.uploaded_amended_orders_id IS 'A foreign key that points to the document table for referencing amended orders';

CREATE INDEX orders_uploaded_amended_orders_id_idx ON orders (uploaded_amended_orders_id);
