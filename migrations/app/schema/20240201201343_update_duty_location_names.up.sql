-- update the duty_location_names values to properly associate with the value on the duty_locations table
UPDATE duty_location_names SET name = 'Ft Novosel' WHERE name = 'Ft Rucker';
UPDATE duty_location_names SET name = 'Ft Gregg-Adams' WHERE name = 'Ft Lee';
UPDATE duty_location_names SET name = 'Ft Cavazos' WHERE name = 'Ft Hood';
UPDATE duty_location_names SET name = 'Ft Liberty' WHERE name = 'Ft Bragg';
UPDATE duty_location_names SET name = 'Ft Johnson' WHERE name = 'Ft Polk';
UPDATE duty_location_names SET name = 'Ft Eisenhower' WHERE name = 'Ft Gordon';
