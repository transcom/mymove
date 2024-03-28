
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
    '5f32e188-9b91-48c9-b696-b02b35ab50ad',
    'db6841f21b5002281512265641bafad5c3f388144ed5f38cf7d4bd3a5491a56b' || '@api.move.mil',
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
    '5f32e188-9b91-48c9-b696-b02b35ab50ad',
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
    'b05b50f6-9825-494f-877c-14d10b32c0dd',
    'db6841f21b5002281512265641bafad5c3f388144ed5f38cf7d4bd3a5491a56b',
    'CN=cameroncaci,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
    '5f32e188-9b91-48c9-b696-b02b35ab50ad',
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
