-- We need access to a UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS service_items_customer_contacts (
	id UUID PRIMARY KEY,
	mtoservice_item_id UUID NOT NULL,
	mtoservice_item_customer_contact_id UUID NOT NULL,
	FOREIGN KEY (mtoservice_item_id) REFERENCES mto_service_items (id),
	FOREIGN KEY (mtoservice_item_customer_contact_id) REFERENCES mto_service_item_customer_contacts (id),
	created_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP without time zone NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP without time zone
);

-- Get all destination SIT service items and any corresponding customer contact ids
-- Some of these will have no customer contacts,
-- but they might have a corresponding destination SIT service item that does have customer contacts
WITH all_dest_sits as (
	SELECT
		si.mto_shipment_id,
		si.id,
		re.code,
		cc.id as cc_id,
		cc.type	as cc_type
	FROM
		mto_service_items si
			INNER JOIN re_services re ON si.re_service_id = re.id
			LEFT join mto_service_item_customer_contacts cc ON si.id = cc.mto_service_item_id
	WHERE re.code IN ('DDASIT', 'DDDSIT', 'DDFSIT')
),
-- In this CTE each row has all destination SIT service items associated with their corresponding shipment
-- and with the first corresponding FIRST and SECOND typed customer contacts
corresponding_contacts_and_service_items as (
 	SELECT
		DISTINCT msi.mto_shipment_id,
		(SELECT all_dest_sits.id FROM all_dest_sits WHERE all_dest_sits.mto_shipment_id = msi.mto_shipment_id AND all_dest_sits.code = 'DDFSIT' LIMIT 1) as ddfsit_id,
		(SELECT all_dest_sits.id FROM all_dest_sits WHERE all_dest_sits.mto_shipment_id = msi.mto_shipment_id AND all_dest_sits.code = 'DDASIT' LIMIT 1) as ddasit_id,
		(SELECT all_dest_sits.id FROM all_dest_sits WHERE all_dest_sits.mto_shipment_id = msi.mto_shipment_id AND all_dest_sits.code = 'DDDSIT' LIMIT 1) as dddsit_id,
		(SELECT all_dest_sits.cc_id FROM all_dest_sits WHERE all_dest_sits.mto_shipment_id = msi.mto_shipment_id AND all_dest_sits.cc_type = 'FIRST' LIMIT 1) as first_cc_id,
		(SELECT all_dest_sits.cc_id FROM all_dest_sits WHERE all_dest_sits.mto_shipment_id = msi.mto_shipment_id AND all_dest_sits.cc_type = 'SECOND' LIMIT 1) as second_cc_id
FROM
    mto_service_items msi
		INNER JOIN re_services re ON msi.re_service_id = re.id
		INNER join mto_service_item_customer_contacts cc ON msi.id = cc.mto_service_item_id
 	WHERE re.code IN ('DDASIT', 'DDDSIT', 'DDFSIT') AND cc.id IS NOT NULL
)
INSERT INTO service_items_customer_contacts (id, mtoservice_item_id, mtoservice_item_customer_contact_id, created_at, updated_at)
-- The FIRST customer contact paired with the DDFSIT
SELECT
	uuid_generate_v4(),
	corresponding_contacts_and_service_items.ddfsit_id   as mtoservice_item_id,
	corresponding_contacts_and_service_items.first_cc_id as mtoservice_item_customer_contact_id,
	NOW(),
	NOW()
FROM corresponding_contacts_and_service_items
WHERE ddfsit_id IS NOT NULL
UNION
-- The SECOND customer contact paired with the DDFSIT
SELECT
	uuid_generate_v4(),
	corresponding_contacts_and_service_items.ddfsit_id    as mtoservice_item_id,
	corresponding_contacts_and_service_items.second_cc_id as mtoservice_item_customer_contact_id,
	NOW(),
	NOW()
FROM corresponding_contacts_and_service_items
WHERE ddfsit_id IS NOT NULL
UNION
-- The FIRST customer contact paired with the DDASIT
SELECT
	uuid_generate_v4(),
	corresponding_contacts_and_service_items.ddasit_id   as mtoservice_item_id,
	corresponding_contacts_and_service_items.first_cc_id as mtoservice_item_customer_contact_id,
	NOW(),
	NOW()
FROM corresponding_contacts_and_service_items
WHERE ddasit_id IS NOT NULL
UNION
-- The SECOND customer contact paired with the DDASIT
SELECT
	uuid_generate_v4(),
	corresponding_contacts_and_service_items.ddasit_id    as mtoservice_item_id,
	corresponding_contacts_and_service_items.second_cc_id as mtoservice_item_customer_contact_id,
	NOW(),
	NOW()
FROM corresponding_contacts_and_service_items
WHERE ddasit_id IS NOT NULL
UNION
-- The FIRST customer contact paired with the DDDSIT
SELECT
	uuid_generate_v4(),
	corresponding_contacts_and_service_items.dddsit_id   as mtoservice_item_id,
	corresponding_contacts_and_service_items.first_cc_id as mtoservice_item_customer_contact_id,
	NOW(),
	NOW()
FROM corresponding_contacts_and_service_items
WHERE dddsit_id IS NOT NULL
UNION
-- The SECOND customer contact paired with the DDDSIT
SELECT
	uuid_generate_v4(),
	corresponding_contacts_and_service_items.dddsit_id    as mtoservice_item_id,
	corresponding_contacts_and_service_items.second_cc_id as mtoservice_item_customer_contact_id,
	NOW(),
	NOW()
FROM corresponding_contacts_and_service_items
WHERE dddsit_id IS NOT NULL;

-- Delete any orphaned records.
DELETE FROM mto_service_item_customer_contacts
	WHERE id NOT IN (SELECT sicc.mtoservice_item_customer_contact_id FROM service_items_customer_contacts sicc);

-- The mto_service_item_id column is not useful now
ALTER TABLE mto_service_item_customer_contacts
	DROP COLUMN mto_service_item_id;
