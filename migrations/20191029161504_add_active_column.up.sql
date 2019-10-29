ALTER TABLE users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE office_users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE admin_users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE dps_users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE users
SET active = NOT deactivated;

UPDATE office_users
SET active = NOT deactivated;

UPDATE admin_users
SET active = NOT deactivated;

UPDATE dps_users
SET active = NOT deactivated;
