DO $$
DECLARE
	client_cert RECORD;
	prime_role_id UUID;
	new_user_id UUID;

BEGIN
	prime_role_id := (SELECT id FROM roles WHERE role_type = 'prime');

	FOR client_cert IN
	SELECT *
	FROM client_certs
		LOOP
			-- Create a user for each client cert
			new_user_id := uuid_generate_v4();
			INSERT INTO users (
				id,
				login_gov_email,
				created_at,
				updated_at
			)
			VALUES (
				new_user_id,
				client_cert.sha256_digest || '@api.move.mil',
				now(),
				now()
			);

			-- Link each client cert to its corresponding user
			UPDATE client_certs
			SET user_id = new_user_id
			WHERE id = client_cert.id;

			-- Create a user role entry for each new user
			INSERT INTO users_roles (id, role_id, user_id, created_at, updated_at)
			VALUES (uuid_generate_v4(), prime_role_id, new_user_id, now(), now());
		END LOOP;
END $$;
