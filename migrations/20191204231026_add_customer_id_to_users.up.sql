ALTER TABLE users
    ADD COLUMN customer_id uuid;
ALTER TABLE users
    ADD CONSTRAINT customer_id_fk FOREIGN KEY (customer_id) REFERENCES customers (id);