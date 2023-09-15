SET LOCAL statement_timeout = 20000;
UPDATE addresses
	SET country = 'US'
	WHERE country = 'United States';

