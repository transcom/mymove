ALTER TABLE mto_shipments
	ADD COLUMN scheduled_delivery_date date,
    add column actual_delivery_date date;

COMMENT ON COLUMN mto_shipments.scheduled_delivery_date IS 'The delivery date the Prime contractor schedules for a shipment after consultation with the customer';
COMMENT ON COLUMN mto_shipments.actual_delivery_date IS 'The actual date that the shipment was delivered to the destination address by the Prime';
