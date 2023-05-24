ALTER TABLE client_certs
	  ADD CONSTRAINT client_certs_subject_idx unique(subject);

ALTER TABLE client_certs
	  ADD CONSTRAINT client_certs_sha256_digest_idx unique(sha256_digest);
