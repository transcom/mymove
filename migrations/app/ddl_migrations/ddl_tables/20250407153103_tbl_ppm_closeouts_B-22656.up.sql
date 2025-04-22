--B-22656  Daniel Jordan  Add small package columns

ALTER TABLE ppm_closeouts
	ADD COLUMN IF NOT EXISTS gtcc_paid_small_package integer,
	ADD COLUMN IF NOT EXISTS member_paid_small_package integer;

COMMENT on COLUMN ppm_closeouts.gtcc_paid_small_package IS 'Amount paid for small package expenses using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_small_package IS 'Amount paid for small package expenses by the service member. Stored in cents.';