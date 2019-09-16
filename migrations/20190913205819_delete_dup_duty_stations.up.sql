-- Minot Air Force Base
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

-- Warren
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

-- Warren Air Force Base
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
