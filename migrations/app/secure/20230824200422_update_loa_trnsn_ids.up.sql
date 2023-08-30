-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

-- Update loa_trnsn_id column constraint
ALTER TABLE lines_of_accounting ALTER COLUMN loa_trnsn_id TYPE varchar (3);
-- Update lines_of_accounting with updated loa_trnsn_id values, mapped by loa_sys_id
UPDATE lines_of_accounting AS loas SET
	loa_trnsn_id = updated.loa_trnsn_id
FROM (VALUES
	(10002, NULL),
	(10001, 'A1'),
	(10003, 'B1'),
	(10004, NULL),
	(10005, 'C1'),
	(10006, 'AB'),
	(10007, NULL),
	(10008, 'C12')
) AS updated(loa_sys_id, loa_trnsn_id)
WHERE updated.loa_sys_id = loas.loa_sys_id;
