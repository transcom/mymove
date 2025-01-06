-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

-- Create test_sequence in migration, moved from sequencer_test.go
CREATE SEQUENCE IF NOT EXISTS test_sequence;