-- We need access to a UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE service_items_customer_contacts (
	id UUID PRIMARY KEY,
	mto_service_item_id UUID NOT NULL,
	mto_service_item_customer_contact_id UUID NOT NULL,
	FOREIGN KEY (mto_service_item_id) REFERENCES mto_service_items (id),
	FOREIGN KEY (mto_service_item_customer_contact_id) REFERENCES mto_service_item_customer_contacts (id),
	created_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP without time zone
);

INSERT INTO service_items_customer_contacts (id, mto_service_item_id, mto_service_item_customer_contact_id)
SELECT uuid_generate_v4(), msicc.mto_service_item_id, msicc.id FROM mto_service_item_customer_contacts msicc;

ALTER TABLE mto_service_item_customer_contacts
	DROP COLUMN mto_service_item_id
