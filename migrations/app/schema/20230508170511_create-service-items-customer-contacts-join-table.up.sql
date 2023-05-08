CREATE TABLE service_items_customer_contacts (
	id UUID PRIMARY KEY,
	mto_service_item_id UUID NOT NULL,
	customer_contact_id UUID NOT NULL,
	FOREIGN KEY (mto_service_item_id) REFERENCES mto_service_items (id),
	FOREIGN KEY (customer_contact_id) REFERENCES mto_service_item_customer_contacts (id),
	created_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP without time zone
);

INSERT INTO service_items_customer_contacts (mto_service_item_id, customer_contact_id)
SELECT msicc.mto_service_item_id, msicc.id FROM mto_service_item_customer_contacts msicc;
