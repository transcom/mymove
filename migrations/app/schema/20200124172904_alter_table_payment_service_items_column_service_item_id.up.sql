ALTER TABLE payment_service_items
    ADD COLUMN mto_service_item_id UUID;

UPDATE payment_service_items SET mto_service_item_id = service_item_id;

ALTER TABLE payment_service_items
    ALTER COLUMN service_item_id DROP NOT NULL;

CREATE INDEX ON payment_service_items (mto_service_item_id);
CREATE INDEX ON payment_service_items (payment_request_id);
