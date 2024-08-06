
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
INSERT INTO users (
    id,
    okta_email,
    created_at,
    updated_at)
VALUES (
    'b64ad6b2-3229-4204-9d21-8031988caf60',
    '04292dab3a650912fa23a339806621b90e2e1da6601180fcd0e33ce27c0cabd9' || '@api.move.mil',
    now(),
    now());

INSERT INTO users_roles (
    id,
    role_id,
    user_id,
    created_at,
    updated_at)
VALUES (
    uuid_generate_v4(),
    (SELECT id FROM roles WHERE role_type = 'prime'),
    'b64ad6b2-3229-4204-9d21-8031988caf60',
    now(),
    now());

INSERT INTO public.client_certs (
    id,
    sha256_digest,
    subject,
    user_id,
    allow_orders_api,
    allow_prime,
    created_at,
    updated_at,
    allow_air_force_orders_read,
    allow_air_force_orders_write,
    allow_army_orders_read,
    allow_army_orders_write,
    allow_coast_guard_orders_read,
    allow_coast_guard_orders_write,
    allow_marine_corps_orders_read,
    allow_marine_corps_orders_write,
    allow_navy_orders_read,
    allow_navy_orders_write)
VALUES (
    '913d7488-fded-424e-9bf7-7f14dde3d596',
    '04292dab3a650912fa23a339806621b90e2e1da6601180fcd0e33ce27c0cabd9',
    'C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=CONTRACTOR, CN=maidocaci',
    'b64ad6b2-3229-4204-9d21-8031988caf60',
    true,
    true,
    now(),
    now(),
    true,
    true,
    true,
    true,
    true,
    true,
    true,
    true,
    true,
    true);
