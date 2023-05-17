ALTER TABLE re_contract_years
	DROP CONSTRAINT re_contract_years_daterange_excl,
	ADD CONSTRAINT re_contract_years_daterange_excl EXCLUDE USING gist(DATERANGE(start_date, end_date, '[]') WITH &&);
