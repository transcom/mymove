CREATE TABLE prime_requesters
(
    id uuid UNIQUE NOT NULL,
    name varchar(255) NOT NULL,
    client_cert_id uuid NOT NULL
        constraint prime_requester_client_cert_id_fkey references client_certs,
    last_seen_at timestamp WITH TIME ZONE NOT NULL,
    created_at timestamp WITH TIME ZONE NOT NULL,
    updated_at timestamp WITH TIME ZONE NOT NULL,
    allow_access bool default FALSE,
    PRIMARY KEY(name, client_cert_id)
);
