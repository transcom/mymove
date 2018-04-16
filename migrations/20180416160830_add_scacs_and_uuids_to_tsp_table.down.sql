DROP TABLE transportation_service_providers;

CREATE TABLE transportation_service_providers (
    id uuid NOT NULL,
    standard_carrier_alpha_code text NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
