ALTER TABLE ppm_shipments
	ADD COLUMN has_office_adjusted_advance bool;

COMMENT ON COLUMN ppm_shipments.has_office_adjusted_advance IS 'An indicator that an office user has denied or adjusted the amount of the requested advance';
