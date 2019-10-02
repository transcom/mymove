-- Minot Air Force Base -> Minot AFB
UPDATE service_members
	SET duty_station_id = (select id from duty_stations where name = 'Minot AFB')
	WHERE duty_station_id = (select id from duty_stations where name = 'Minot Air Force Base');
UPDATE orders
	SET new_duty_station_id = (select id from duty_stations where name = 'Minot AFB')
	WHERE new_duty_station_id = (select id from duty_stations where name = 'Minot Air Force Base');
DELETE FROM duty_stations WHERE name = 'Minot Air Force Base'
	AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = (select id from duty_stations where name = 'Minot Air Force Base'))
	AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = (select id from duty_stations where name = 'Minot Air Force Base'));
DELETE FROM addresses WHERE id = '43e3ab4a-a307-47da-b2dc-93a1bc8ace44'
	AND NOT EXISTS (SELECT id FROM duty_stations WHERE name = 'Minot Air Force Base');

-- Warren -> F.E. Warren AFB
UPDATE service_members
	SET duty_station_id = (select id from duty_stations where name = 'F.E. Warren AFB')
	WHERE duty_station_id = (select id from duty_stations where name = 'Warren');
UPDATE orders
	SET new_duty_station_id = (select id from duty_stations where name = 'F.E. Warren AFB')
	WHERE new_duty_station_id = (select id from duty_stations where name = 'Warren');
DELETE FROM duty_stations WHERE name = 'Warren'
	AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = (select id from duty_stations where name = 'Warren'))
	AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = (select id from duty_stations where name = 'Warren'));
DELETE FROM addresses WHERE id = '3ff01b01-a1ce-4442-b588-a434858bf8c5'
	AND NOT EXISTS (SELECT id FROM duty_stations WHERE name = 'Warren');

-- Warren Air Force Base -> F.E. Warren AFB
UPDATE service_members
	SET duty_station_id = (select id from duty_stations where name = 'F.E. Warren AFB')
	WHERE duty_station_id = (select id from duty_stations where name = 'Warren Air Force Base');
UPDATE orders
	SET new_duty_station_id = (select id from duty_stations where name = 'F.E. Warren AFB')
	WHERE new_duty_station_id = (select id from duty_stations where name = 'Warren Air Force Base');
DELETE FROM duty_stations WHERE name = 'Warren Air Force Base'
	AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = (select id from duty_stations where name = 'Warren Air Force Base'))
	AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = (select id from duty_stations where name = 'Warren Air Force Base'));
DELETE FROM addresses WHERE id = '44c1df72-c433-4d09-a8c0-89672fc2f3ee'
	AND NOT EXISTS (SELECT id FROM duty_stations WHERE name = 'Warren Air Force Base');

-- USCG Mobile -> Sector Mobile
UPDATE service_members
	SET duty_station_id = (select id from duty_stations where name = 'Sector Mobile')
	WHERE duty_station_id = (select id from duty_stations where name = 'USCG Mobile');
UPDATE orders
	SET new_duty_station_id = (select id from duty_stations where name = 'Sector Mobile')
	WHERE new_duty_station_id = (select id from duty_stations where name = 'USCG Mobile');
DELETE FROM duty_stations WHERE name = 'USCG Mobile'
	AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = (select id from duty_stations where name = 'USCG Mobile'))
	AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = (select id from duty_stations where name = 'USCG Mobile'));
DELETE FROM addresses WHERE id = 'be84b14a-f836-4821-af13-fdbe8530386d'
	AND NOT EXISTS (SELECT id FROM duty_stations WHERE name = 'USCG Mobile');

-- Fort Eustis -> JB Langley-Eustis
UPDATE service_members
	SET duty_station_id = (select id from duty_stations where name = 'JB Langley-Eustis')
	WHERE duty_station_id = (select id from duty_stations where name = 'Fort Eustis');
UPDATE orders
	SET new_duty_station_id = (select id from duty_stations where name = 'JB Langley-Eustis')
	WHERE new_duty_station_id = (select id from duty_stations where name = 'Fort Eustis');
DELETE FROM duty_stations WHERE name = 'Fort Eustis'
	AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = (select id from duty_stations where name = 'Fort Eustis'))
	AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = (select id from duty_stations where name = 'Fort Eustis'));
DELETE FROM addresses WHERE id = '4c1605b8-ac15-41ef-bf1e-b26d45c0005e'
	AND NOT EXISTS (SELECT id FROM duty_stations WHERE name = 'Fort Eustis');