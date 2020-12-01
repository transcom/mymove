ALTER TABLE mto_service_items
	ADD COLUMN approved_at timestamp without time zone,
    ADD COLUMN rejected_at timestamp without time zone;