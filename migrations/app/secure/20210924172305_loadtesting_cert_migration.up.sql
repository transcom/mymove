
-- This migration adds a cert that can access the MilMove Loadtest environment with mTLS.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
INSERT INTO public.client_certs (
	id,
	sha256_digest,
	subject,
	allow_dps_auth_api,
	allow_orders_api,
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
	allow_navy_orders_write,
	allow_prime)
VALUES (
	'debb6c7f-bfdf-4621-ade4-8def4571034a',
	'215732d45860c0d01a00a8289ea57bc2391f8ed7097f8adbfc924c5b1b47c517',
	'CN=api.loadtest.dp3.us',
	false,
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
	true,
	true);
