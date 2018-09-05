-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO public.office_users
VALUES
    ('cc90e5ef-601c-4e6c-af82-15b8cdc3ab90', NULL, 'Ripley', 'Ellen', NULL, 'ripley@nostromo.space', '(555) 555-5555', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users
VALUES
    ('8293e5d2-819c-4597-956e-b43072f5c488', NULL, 'Android', 'Ash', NULL, 'ash@hyperdyne.biz', '(555) 555-5555', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());