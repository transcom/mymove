-- PRD ATO migration
-- This will be distributed ATO PRD environments only

-- DO $$
-- DECLARE
-- 	daycos_user_id UUID;
-- 	daycos_client_cert_id UUID;
-- 	daycos_users_roles_id UUID;

-- 	movehq_user_id UUID;
-- 	movehq_client_cert_id UUID;
-- 	movehq_users_roles_id UUID;

-- BEGIN
-- 	daycos_user_id := <UUID>;
-- 	daycos_client_cert_id := <UUID>;
-- 	daycos_users_roles_id := <UUID>;

-- 	movehq_user_id := <UUID>;
-- 	movehq_client_cert_id := <UUID>;
-- 	movehq_users_roles_id := <UUID>;

-- 	DELETE FROM users WHERE id = daycos_user_id LIMIT 1;
-- 	DELETE FROM users WHERE id = movehq_user_id LIMIT 1;
-- 	DELETE FROM client_certs WHERE id = daycos_client_cert_id LIMIT 1;
-- 	DELETE FROM client_certs WHERE id = movehq_client_cert_id LIMIT 1;
-- 	DELETE FROM users_roles WHERE id = daycos_users_roles_id LIMIT 1;
-- 	DELETE FROM users_roles WHERE id = movehq_users_roles_id LIMIT 1;

-- END $$;