
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
DELETE FROM client_certs WHERE user_id = '545f3a07-76dd-4c62-8541-0353cb507d17';
DELETE FROM users_roles WHERE user_id = '545f3a07-76dd-4c62-8541-0353cb507d17';
DELETE FROM users WHERE id = '545f3a07-76dd-4c62-8541-0353cb507d17';

INSERT INTO users (
    id,
    okta_email,
    created_at,
    updated_at)
VALUES (
    '545f3a07-76dd-4c62-8541-0353cb507d17',
    'dd28f2ed02b4ed5065e7d72817373303c8a2de424c1902c1c5afe16309956a56' || '@api.move.mil',
    now(),
    now());

INSERT INTO users_roles (
    id,
    role_id,
    user_id,
    created_at,
    updated_at)
VALUES (
    '545f3a07-76dd-4c62-8541-0353cb507d17',
    (SELECT id FROM roles WHERE role_type = 'prime'),
    '545f3a07-76dd-4c62-8541-0353cb507d17',
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
    'cf2b7400-0890-400d-85a2-906ce34281f3',
    'dd28f2ed02b4ed5065e7d72817373303c8a2de424c1902c1c5afe16309956a56',
    'CN=joeydoyecaci,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
    '545f3a07-76dd-4c62-8541-0353cb507d17',
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
