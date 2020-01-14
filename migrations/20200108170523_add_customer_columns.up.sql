ALTER TABLE customers
    ADD COLUMN email                  text,
    ADD COLUMN phone                  text,
    ADD COLUMN current_address_id     uuid,
    ADD COLUMN destination_address_id uuid;

ALTER TABLE customers
    ADD CONSTRAINT customers_current_address_fk FOREIGN KEY (current_address_id) REFERENCES addresses (id),
    ADD CONSTRAINT customers_destination_address_fk FOREIGN KEY (destination_address_id) REFERENCES addresses (id);