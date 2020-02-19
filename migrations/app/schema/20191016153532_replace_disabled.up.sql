ALTER TABLE users
ADD COLUMN deactivated BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE office_users
ADD COLUMN deactivated BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE admin_users
ADD COLUMN deactivated BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE dps_users
ADD COLUMN deactivated BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE users
SET deactivated = disabled;

UPDATE office_users
SET deactivated = disabled;

UPDATE admin_users
SET deactivated = disabled;

UPDATE dps_users
SET deactivated = disabled;
