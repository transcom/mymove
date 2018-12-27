-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO public.client_certs VALUES ('cb8ed5ed41787e84c5f01dc8599fb83b8397cdc0a5eda8f532eae52b5682c342', 'C=US, ST=DC, L=Washington, O=Faux Air Force, OU=Faux Orders, CN=localhost', false, true, now(), now());
