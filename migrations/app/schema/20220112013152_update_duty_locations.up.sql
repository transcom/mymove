-- There were a few more bases that needed to be deleted that didn't make it into the previous migration
-- Slack thread for context: https://ustcdp3.slack.com/archives/C02AXHB7YRE/p1641502874000100?thread_ts=1640902686.117800&cid=C02AXHB7YRE
DROP TABLE IF EXISTS list_delete;
CREATE TEMP TABLE list_delete
(id uuid)
	ON COMMIT DROP;

INSERT INTO list_delete VALUES
							-- 33D FIGHTER WING HQTRS
							('d19b613c-2a40-46bc-9ee1-f229db255eb2'),
							-- CABLE SPLICER INST SHEPPARD AFB TX
							('ccce11b1-c758-4b8c-8ed9-aaf779d1a99e'),
							-- DET MARINE AVIATION NAS PATUXENT RIVER MD
							('cef7572a-5771-428d-bacb-32f20e78d475'),
							-- MARCOR DET FT LEAVENWORTH KS
							('82fcdd56-5527-44ea-b8d2-bb90a1953384'),
							-- USN POST GRADUATE SCHOOL (STUD PERS)
							('f9fc7e21-fa66-4b1e-b26a-4efcbb7bf22c');

-- Change duty location to Los Angeles AFB for any service members associated with locations that will be deleted.
UPDATE service_members
SET duty_station_id = 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96'
WHERE duty_station_id IN (select id from list_delete);

-- For any orders that have either origin or destination duty location set to a location that will be deleted,
-- change their origin_duty_station to Los Angeles AFB (GBLOC KKFA) and their new_duty_station to Fort Bragg.
-- The goal is to not break existing moves too badly, without putting too much effort into maintaining them.
UPDATE orders
SET new_duty_station_id = 'dca78766-e76b-4c6d-ba82-81b50ca824b9', origin_duty_station_id = 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96'
WHERE new_duty_station_id IN (select id from list_delete) OR origin_duty_station_id IN (select id from list_delete);

DELETE FROM duty_station_names WHERE duty_station_id IN (select id from list_delete);
DELETE FROM duty_locations WHERE id IN (select id from list_delete);
