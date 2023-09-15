-- -- STG ATO migration
-- -- This will be distributed ATO STG environments only

-- DO $$
-- DECLARE
-- 	movehq_user_id UUID;
-- 	new_movehq_sha256 VARCHAR;
-- 	movehq_client_cert_id UUID;

-- BEGIN

-- 	movehq_user_id := '<UUID>';
-- 	new_movehq_sha256 := '<SHA256>';
-- 	movehq_client_cert_id := '<UUID>';

-- 	UPDATE users
-- 		SET
-- 			login_gov_email = new_movehq_sha256 || '@api.move.mil'
-- 		WHERE id = movehq_user_id;

-- 	UPDATE client_certs
-- 		SET
-- 			sha256_digest = new_movehq_sha256,
-- 			subject = 'subject= CN=mmb.gov.uat.homesafeconnect.com,OU=MoveHQ Inc.,OU=IdenTrust,OU=ECA,O=U.S. Government,C=US'
-- 		WHERE id = movehq_client_cert_id;

-- END $$;
