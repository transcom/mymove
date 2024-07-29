
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
    '458c3c7b-a047-4c79-b68b-552e6e7bba32',
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
    uuid_generate_v4(),
    (SELECT id FROM roles WHERE role_type = 'prime'),
    '458c3c7b-a047-4c79-b68b-552e6e7bba32',
    now(),
    now());

INSERT INTO public.client_certs (
    id,
    sha256_digest,
    subject,
    user_id,
    allow_orders_api,
    allow_prime,
	allow_pptas,
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
    'c00707e0-98d0-47a3-a65a-00f2ddfde60f',
    'dd28f2ed02b4ed5065e7d72817373303c8a2de424c1902c1c5afe16309956a56',
    'CN=joeydoyecaci2,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
    '458c3c7b-a047-4c79-b68b-552e6e7bba32',
    true,
    true,
	false,
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
