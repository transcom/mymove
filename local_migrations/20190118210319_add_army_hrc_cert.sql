-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO public.client_certs VALUES ('c9f663b7-3c37-4c9b-b1b8-61e3b7155e5f', '082f91f927e78847b8f74010427bd00ac9e3313e9f29bff6f637dea91c815f12', '/C=US/ST=DC/L=Washington/O=Not Army HRC/OU=Not Army HRC Orders/CN=localhost', false, true, now(), now());
