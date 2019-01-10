-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO public.client_certs VALUES ('f885b0d3-3df4-46b3-908e-c9c3fec9d2f4', '493ba2a4634b002d3f093e88bd182ce885e04d7efa6132b1fcfbb14055bf66e6', '/C=US/ST=DC/L=Washington/O=Faux Marine Corps/OU=Faux Orders/CN=localhost', false, true, now(), now());
