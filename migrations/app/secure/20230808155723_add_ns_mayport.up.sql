--- This migration is only expected to run in the STG environment. The other environments have a blank file.

UPDATE addresses SET postal_code = '32233' WHERE id = '430c50b4-cc87-42a6-9921-2010aa6d4ed3';
INSERT INTO duty_locations (id, address_id, name, affiliation, transportation_office_id, updated_at, created_at, provides_services_counseling) VALUES ('3bfe8068-f9d1-4c21-ac19-1c0a60dcd326', '430c50b4-cc87-42a6-9921-2010aa6d4ed3', 'NS Mayport, FL 32233', 'NAVY', 'f5ab88fe-47f8-4b58-99af-41067d6cb60d', now(), now(), 'TRUE');
