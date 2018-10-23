-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

-- Enroll TSPs in the HHG program and add data
-- "supplier_id" fields do not match production fields

UPDATE transportation_service_providers
SET
	enrolled = true,
	name = 'Planetary Van Lines, LLC',
	supplier_id = 'PYVL1234',
	poc_general_name = 'Joey Jupiter',
	poc_general_email = 'j.jupiter@example.com',
	poc_general_phone = '(555) 123-4567',
	poc_claims_name = 'Chaim Calypso',
	poc_claims_email = 'claims@example.com',
	poc_claims_phone = '(555) 765-4321'
WHERE standard_carrier_alpha_code = 'PYVL';

UPDATE transportation_service_providers
SET
	enrolled = true,
	name = 'Green Chip LLC dba D’Lux Moving and Storage',
	supplier_id = 'DLXM1234',
	poc_general_name = 'Gary Green',
	poc_general_email = 'g.green@example.com',
	poc_general_phone = '(555) 124-8163',
	poc_claims_name = 'Cassandra Gregory',
	poc_claims_email = 'claims@example.com',
	poc_claims_phone = '(555) 264-1282'
WHERE standard_carrier_alpha_code = 'DLXM';

UPDATE transportation_service_providers
SET
	enrolled = true,
	name = 'Secure Storage Company of Washington, LLC',
	supplier_id = 'SSOW1234',
	poc_general_name = 'Seth Stallone',
	poc_general_email = 's.stallone@example.com',
	poc_general_phone = '(555) 101-0101',
	poc_claims_name = 'Charlotte Salsbury',
	poc_claims_email = 'claims@example.com',
	poc_claims_phone = '(555) 111-0000'
WHERE standard_carrier_alpha_code = 'SSOW';


