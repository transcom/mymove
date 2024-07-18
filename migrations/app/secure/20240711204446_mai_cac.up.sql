-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prd/stg/exp/demo
-- DO NOT include any sensitive data.

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
    '4cb57e6e-e635-4700-8e46-08338452d56b',
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
    '4cb57e6e-e635-4700-8e46-08338452d56b',
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
    '468a2081-de9d-49cc-8a12-5c661ce53ca1',
    '04292dab3a650912fa23a339806621b90e2e1da6601180fcd0e33ce27c0cabd9',
    'C=US, O=U.S. Government, OU=DoD, OU=PKI, OU=CONTRACTOR, CN=maidocaci',
    '4cb57e6e-e635-4700-8e46-08338452d56b',
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
