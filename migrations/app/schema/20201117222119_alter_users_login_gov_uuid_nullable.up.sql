-- this was commented out due to removal of the login_gov columns in an earlier migration file
ALTER TABLE users ALTER COLUMN login_gov_uuid DROP NOT NULL;
ALTER TABLE users ALTER COLUMN login_gov_email DROP NOT NULL;