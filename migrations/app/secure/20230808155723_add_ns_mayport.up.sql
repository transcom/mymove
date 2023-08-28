--- This migration is only expected to run in the STG environment. The other environments have a blank file.

UPDATE addresses SET postal_code = '32233' WHERE id = '430c50b4-cc87-42a6-9921-2010aa6d4ed3';
UPDATE duty_locations SET name = 'NS Mayport, FL 32233', provides_services_counseling = 'TRUE' where address_id = '430c50b4-cc87-42a6-9921-2010aa6d4ed3'
