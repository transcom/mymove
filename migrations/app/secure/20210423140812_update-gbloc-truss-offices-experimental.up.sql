-- This is a migration that will ultimately impact experimental. It seeks to change the gbloc of a commonly used
-- test transportation office.
UPDATE transportation_offices
SET gbloc = 'AGFM'
WHERE name = 'Truss'
