CREATE TYPE service_item_status AS ENUM (
	'SUBMITTED',
	'APPROVED',
	'REJECTED'
);

ALTER TABLE mto_service_items ADD COLUMN status service_item_status NOT NULL DEFAULT 'SUBMITTED';
