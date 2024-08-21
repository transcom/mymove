ALTER TABLE admin_users
ADD COLUMN IF NOT EXISTS super BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT on COLUMN admin_users.super IS 'Value that designates super admin users.';