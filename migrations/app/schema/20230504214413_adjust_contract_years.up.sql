ALTER TABLE re_contract_years
	DROP CONSTRAINT re_contract_years_daterange_excl,
	ADD CONSTRAINT re_contract_years_daterange_excl EXCLUDE USING gist(DATERANGE(start_date, end_date, '[]') WITH &&);

-- Update contract years for TRUSS_TEST so that they won't overlap new records
UPDATE re_contract_years
SET end_date = '2023-04-01'
WHERE id = '9ac23b92-742a-4d17-9a50-61ae4cf3f3e3';

UPDATE re_contract_years
SET start_date = '2023-04-02',
	end_date   = '2023-04-03'
WHERE id = 'a6457124-b8d5-4b75-95bd-fd402974a043';

UPDATE re_contract_years
SET start_date = '2023-04-04',
	end_date   = '2023-04-05'
WHERE id = '2546c4de-25dd-499d-a6da-61e4de247649'

UPDATE re_contract_years
SET start_date = '2023-04-06',
	end_date   = '2023-04-07'
WHERE id = '53b69fab-c9d6-4900-9998-03b4980519ba';

-- We should update this to be contiguous with the Test Period for the new pricing data when we have that
UPDATE re_contract_years
SET start_date = '2023-04-08',
	end_date   = '2023-05-10'
WHERE id = '741f66ee-34c6-4388-b6ce-c48b6597b6e3';
