-- This migration creates records in the users_roles table for each office user
-- that currently doesn't have an explicit role, and assumes that the default role_id :=
-- for existing users is 'ppm_office_users'

DO $$
DECLARE
	rec RECORD;
	role_id UUID;

BEGIN
	role_id := (SELECT id FROM roles WHERE role_type = 'ppm_office_users');
	-- Find users that have associated office user records that have no explicit role
	FOR rec IN SELECT DISTINCT users.*
		FROM users
		INNER JOIN office_users ON users.id = office_users.user_id
		LEFT JOIN users_roles ON users.id = users_roles.user_id
		WHERE users_roles.user_id IS NULL

		LOOP
			INSERT INTO users_roles (id, role_id, user_id, created_at, updated_at)
			VALUES (uuid_generate_v4(), role_id, rec.id, now(), now());
	END LOOP;
END $$;
