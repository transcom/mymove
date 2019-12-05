CREATE TYPE payment_service_item_status AS ENUM (
    'REQUESTED',
    'APPROVED',
    'DENIED',
    'SENT_TO_GEX',
    'PAID'
    );

CREATE TABLE payment_service_items
(
    id uuid primary key,
    payment_request_id uuid NOT NULL
        constraint payment_service_items_payment_request_id_fkey references payment_requests,
    service_item_id uuid,
    status payment_service_item_status NOT NULL,
    price_cents integer NOT NULL,
    reject_reason varchar(255),
    requested_at timestamp without time zone NOT NULL,
    approved_at timestamp without time zone,
    denied_at timestamp without time zone,
    sent_to_gex_at timestamp without time zone,
    paid_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
)