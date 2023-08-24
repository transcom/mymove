ALTER TABLE mto_service_items
	ADD COLUMN sit_customer_contacted date,
	ADD COLUMN sit_requested_delivery date,

-- Comment On Column
COMMENT ON COLUMN mto_service_items.sit_customer_contacted IS 'The date when the customer contacted the prime for a delivery out of SIT';
COMMENT ON COLUMN mto_service_items.sit_requested_delivery IS 'The date when the customer has requested delivery out of SIT';
