ALTER TABLE payment_service_items
    ALTER COLUMN service_item_id SET NOT NULL,
    ADD CONSTRAINT payment_service_items_service_item_id_fkey FOREIGN KEY (service_item_id) REFERENCES mto_service_items (id);
