ALTER TABLE users
DROP COLUMN disabled;

ALTER TABLE office_users
DROP COLUMN disabled;

ALTER TABLE admin_users
DROP COLUMN disabled;

ALTER TABLE dps_users
DROP COLUMN disabled;
