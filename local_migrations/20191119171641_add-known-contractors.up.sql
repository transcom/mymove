-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.
INSERT INTO contractor VALUES ('5db13bb4-6d29-4bdb-bc81-262f4513ecf6',  now(), now(), 'Prime McPrime Contractor', 'HTC111-11-1-1111','Prime');
INSERT INTO contractor VALUES ('ee32183b-7bc2-4587-97c4-9863c2a6937a',  now(), now(), 'NTS Contractor 1', 'NTC111-11-1-1111','NTS');