UPDATE duty_stations SET affiliation='AIR_FORCE' WHERE affiliation='AIRFORCE';
UPDATE duty_stations SET affiliation='COAST_GUARD' WHERE affiliation='COASTGUARD';

UPDATE transportation_offices SET name='Kirtland AFB' WHERE name='Kirkland AFB';

-- City was set as "Kortland AFB"
UPDATE addresses SET city='Kirtland AFB' WHERE id=(SELECT address_id FROM duty_stations WHERE name='Kirtland AFB');

-- This association was never set up though the transportation offices existed
UPDATE duty_stations SET transportation_office_id=(SELECT id FROM transportation_offices WHERE name='Kirtland AFB') WHERE name='Kirtland AFB';
