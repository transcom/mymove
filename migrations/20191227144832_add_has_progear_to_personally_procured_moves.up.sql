CREATE TYPE progear_status AS ENUM (
    'YES',
    'NO',
    'NOT SURE'
    );

ALTER TABLE personally_procured_moves
	ADD COLUMN has_pro_gear progear_status,
	ADD COLUMN has_pro_gear_over_thousand progear_status;