-- Add PAYMENT_REQUEST origin for a service item param key. Note that this has to be in a
-- separate migration/transaction and applied/committed before it can be used.
ALTER TYPE service_item_param_origin
	ADD VALUE 'PAYMENT_REQUEST';
