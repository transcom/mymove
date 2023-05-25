-- We need access to a UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS service_items_customer_contacts (
	id UUID PRIMARY KEY,
	mtoservice_item_id UUID NOT NULL,
	mtoservice_item_customer_contact_id UUID NOT NULL,
	FOREIGN KEY (mtoservice_item_id) REFERENCES mto_service_items (id),
	FOREIGN KEY (mtoservice_item_customer_contact_id) REFERENCES mto_service_item_customer_contacts (id),
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITHOUT TIME ZONE
);

-- Get one pair of shipments and Customer Contacts per contact type
WITH shipment_contact_pairs AS (
	SELECT DISTINCT ON ((msi.mto_shipment_id, msicc.type))
		msi.mto_shipment_id,
		msicc.type,
		msicc.id AS contact_id
	FROM
		mto_service_items msi
			INNER JOIN re_services res ON msi.re_service_id = res.id
			LEFT JOIN mto_service_item_customer_contacts msicc ON msi.id = msicc.mto_service_item_id
	WHERE res.code IN ('DDASIT', 'DDDSIT', 'DDFSIT') AND msicc.id IS NOT NULL
	ORDER BY (msi.mto_shipment_id, msicc.type)
)
INSERT INTO service_items_customer_contacts (id, mtoservice_item_id, mtoservice_item_customer_contact_id, created_at, updated_at)
SELECT
	uuid_generate_v4(),
	msi.id,
	shipment_contact_pairs.contact_id,
	NOW(),
	NOW()
-- Get each destination service item associated with the shipments from our CTE and create a record with the corresponding customer contacts
FROM
	mto_service_items msi
		INNER JOIN re_services res ON msi.re_service_id = res.id
		INNER JOIN shipment_contact_pairs on msi.mto_shipment_id = shipment_contact_pairs.mto_shipment_id
WHERE res.code IN ('DDASIT', 'DDDSIT', 'DDFSIT');

-- Delete any orphaned records.
DELETE FROM mto_service_item_customer_contacts
	WHERE id NOT IN (SELECT sicc.mtoservice_item_customer_contact_id FROM service_items_customer_contacts sicc);

-- The mto_service_item_id column is not useful now
ALTER TABLE mto_service_item_customer_contacts
	DROP COLUMN mto_service_item_id;
