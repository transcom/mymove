-- B-17750
-- This will only be run on staging environments
UPDATE duty_locations
	SET provides_services_counseling = 'FALSE'
	WHERE name = 'NS Mayport, FL 32228';