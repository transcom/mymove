-- This migration allows Load Testing to have read/write access to all orders and the Prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Load Testing is only permitted to access deployments in experimental and staging - these
-- permissions should NEVER be applied in production.
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
	'730bf94b-f54e-46c3-b125-0a572d885209',
	'4bd850a6aa6b6f403c63d876b931afa76159e7eb23bb2ba0d19ee6a937df6240',
	'CN=orders.exp.move.mil,OU=USTRANSCOM,OU=PKI,OU=DoD,O=U.S. Government,C=US',
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
