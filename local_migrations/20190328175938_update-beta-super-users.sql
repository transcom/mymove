-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.
UPDATE users SET is_superuser = TRUE, updated_at = now() WHERE login_gov_email = 'officeuser1@example.com';