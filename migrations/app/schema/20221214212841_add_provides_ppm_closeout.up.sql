ALTER TABLE transportation_offices
ADD COLUMN provides_ppm_closeout boolean DEFAULT false NOT NULL;

-- Column comment
COMMENT ON COLUMN transportation_offices.provides_ppm_closeout IS 'Indicates whether a transportation office provides ppm closeout or not. It is used by Army and Air Force service members';
