ALTER TABLE users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE office_users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE admin_users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE dps_users
ADD COLUMN active BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE users
SET active = NOT deactivated;

UPDATE office_users
SET active = NOT deactivated;

UPDATE admin_users
SET active = NOT deactivated;

UPDATE dps_users
SET active = NOT deactivated;

CREATE INDEX users_active_idx ON public.users USING btree (active);
CREATE INDEX office_users_active_idx ON public.office_users USING btree (active);
CREATE INDEX admin_users_active_idx ON public.admin_users USING btree (active);
CREATE INDEX dps_users_active_idx ON public.dps_users USING btree (active);
