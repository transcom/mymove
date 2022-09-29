WITH customer_role_id AS (
    SELECT roles.id FROM roles
    WHERE roles.role_type = 'customer'
)
-- Query for any existing service members that don't have customer roles assigned on the user_roles table.
-- Add a customer role to the users_roles table for those service members.
-- Note: We chose to use uuid_generate_v4 since this id should never be referenced and we have an unknown number of
--       service member records to populate on the users_roles table.
INSERT INTO users_roles (user_id, role_id, id, created_at, updated_at)
	SELECT service_members.user_id, (SELECT * FROM customer_role_id), uuid_generate_v4() as id, now(), now()
	FROM service_members
	WHERE NOT EXISTS(
		SELECT users_roles.user_id FROM users_roles
		WHERE users_roles.role_id = (SELECT * FROM customer_role_id)
			AND users_roles.user_id = service_members.user_id
		)
