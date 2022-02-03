-- Remove monfresh's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest='061a7533ab529e259b7afacd8ff69a3a06d57d2b3b831db924aa8b339415f2b3';
-- Remove shkeating's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest='b3c44ca81bcabee2c7dcdc8c3c6809400468d2c25dc3556e4a0e2e143fc5aa85';
-- Remove ronolibert's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest='af6f3129b664216989e983cf09f723a7617bf9d21333fd66294551cea78ba45f';
-- Remove carterjones's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest='27c54b91080a7252215c68d8808de4a2d39cbc3e73fbecb9bddc621f5dde08f7';
