ALTER TABLE move_orders
    ADD COLUMN confirmation_number text,
    ADD COLUMN order_number text,
    ADD COLUMN grade text;

ALTER TABLE customers
    ADD COLUMN first_name text,
    ADD COLUMN last_name text,
    ADD COLUMN agency text;
