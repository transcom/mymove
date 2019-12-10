CREATE TYPE payment_request_status AS ENUM (
    'PENDING',
    'REVIEWED',
    'SENT_TO_GEX',
    'RECEIVED_BY_GEX',
    'PAID'
    );

ALTER TABLE payment_requests
	ADD COLUMN move_task_order_id uuid
	ADD COLUMN status payment_request_status NOT NULL default 'PENDING',
	ADD COLUMN requested_at timestamp without time zone NOT NULL default Now(),
	ADD COLUMN reviewed_at timestamp without time zone,
	ADD COLUMN sent_to_gex_at timestamp without time zone,
	ADD COLUMN received_by_gex_at timestamp without time zone,
	ADD COLUMN paid_at timestamp without time zone;

ALTER TABLE payment_requests ALTER COLUMN is_final SET NOT NULL;
ALTER TABLE payment_requests ALTER COLUMN rejection_reason TYPE varchar(255);