ALTER TABLE customers
    ADD COLUMN user_id uuid;
ALTER TABLE customers
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE customers
    ADD CONSTRAINT customers_user_id_ukey UNIQUE (user_id);