ALTER TABLE payment_service_items
    DROP CONSTRAINT payment_service_items_service_item_id_fkey,
    DROP COLUMN service_item_id,
    ALTER COLUMN mto_service_item_id SET NOT NULL,
    ADD CONSTRAINT payment_service_items_mto_service_item_id_fkey FOREIGN KEY (mto_service_item_id) REFERENCES mto_service_items (id);