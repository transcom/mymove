UPDATE transportation_offices
    SET provides_ppm_closeout = FALSE
    WHERE name in (
	'Camp LeJeune (USMC)'
);
