
-- Until the admin UI is in place and has visibility on the electronic orders table,
-- we need certificates that can look at the Orders that have been uploaded.
-- This migration allows a CAC cert to have read/write access to all orders.
-- The Orders API uses client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to upload Orders and review Orders. Eventually
-- this CAC certificate should be removed.
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
	allow_navy_orders_write)
VALUES (
	'9ca7b372-3dc9-46f9-8010-e66889420544',
	'bc95ac7064f4d1df17eaf446a8c0569786918e74423636db38dfd05191ef481b',
	'/O=Snake Oil/CN=localhost',
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
	true);
