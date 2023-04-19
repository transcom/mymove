-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.
UPDATE contractors
	SET
		name = 'Test Prime',
		contract_number = 'HTC111-11-1-1112'
	WHERE "type" = 'Prime'
