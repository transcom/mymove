-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO public.client_certs VALUES ('ff001674-6a19-4f59-9417-fcfdfe071272', '320c9dc085725aaa925ad1ab00261dc393264a78705bd6d5fc1c37bc33f285dd', '/C=US/ST=IL/L=Belleville/O=Not USTRANSCOM/OU=Not Defense Personal Property System/CN=localhost', true, false, now(), now());
