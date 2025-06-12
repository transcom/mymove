-- B-22484 Paul Stonebraker PPPO JBER Travel Center Elmendorf and PPPO JBER Travel Center Richardson
-- to be combined into one record to be called: PPPO JBER Elmendorf-Richardson  

-- update the FKs
UPDATE office_phone_lines
SET transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

UPDATE office_emails
SET transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

UPDATE moves
SET closeout_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE closeout_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

UPDATE moves
SET counseling_transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE counseling_transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

UPDATE office_users
SET transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

UPDATE duty_locations
SET transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

UPDATE transportation_office_assignments
SET transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14'
WHERE transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

-- combine the two
UPDATE transportation_offices
SET "name" = 'PPPO JBER Elmendorf-Richardson'
WHERE id = '4522d141-87f1-4f1e-a111-466303c6ae14';

-- delete the duplicate
DELETE FROM transportation_offices
WHERE id = 'bc34e876-7f18-4401-ab91-507b0861a947';