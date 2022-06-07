ALTER TABLE
	client_certs
ADD COLUMN
	user_id uuid CONSTRAINT client_certs_user_id_fkey REFERENCES users (id);

COMMENT ON COLUMN client_certs.user_id IS 'Associate a user with each client cert; initially designed to identify a prime "user" from http requests';
