-- This is a migration that will ultimately impact experimental. It seeks to change the gbloc of a commonly used
-- test transportation office.
UPDATE transportation_offices
SET gbloc = 'AGFM'
WHERE id = 'ffee171d-f085-4451-95b7-918dc4d5fff7'
