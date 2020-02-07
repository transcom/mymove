-- Allows us to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION insert_duty_station(
	name varchar,
	branch varchar,
	city varchar,
	state varchar,
	postal_code varchar)
RETURNS void language plpgsql AS $$
DECLARE
    inserted_id uuid;
BEGIN
    INSERT INTO addresses(id, created_at, updated_at, street_address_1, city, state, postal_code)
    VALUES (uuid_generate_v4(), now(), now(), 'n/a', $3, $4, $5)
    RETURNING id INTO inserted_id;

    INSERT INTO duty_stations(id, created_at, updated_at, name, branch, address_id)
    VALUES (uuid_generate_v4(), now(), now(), $1, $2, inserted_id);
END $$;

-- Data from from this spreadsheet: https://docs.google.com/spreadsheets/d/1hlXB1ZXyQVA3uHGqRkCNojzT-jcjfj-PcAeyeF_vqEw/edit#gid=1342927694
SELECT insert_duty_station('Altus AFB','AIRFORCE','Altus AFB','OK','73523');
SELECT insert_duty_station('Barksdale AFB','AIRFORCE','Barksdale AFB','LA','71110');
SELECT insert_duty_station('Columbus AFB','AIRFORCE','Columbus','MS','39710');
SELECT insert_duty_station('Dyess AFB','AIRFORCE','Dyess AFB','TX','79607');
SELECT insert_duty_station('Eglin AFB','AIRFORCE','Eglin AFB','FL','32542');
SELECT insert_duty_station('Ellington Field ANGB','AIRFORCE','Houston','TX','77034');
SELECT insert_duty_station('Fort Sam Houston','ARMY','San Antonio','TX','78234');
SELECT insert_duty_station('Goodfellow AFB','AIRFORCE','Goodfellow AFB','TX','76908');
SELECT insert_duty_station('Hurlburt Field AFB','AIRFORCE','Hurlburt Field','FL','32544');
SELECT insert_duty_station('Keesler AFB','AIRFORCE','Biloxi','MS','39534');
SELECT insert_duty_station('Lackland AFB','AIRFORCE','Lackland AFB','TX','78236');
SELECT insert_duty_station('Laughlin AFB','AIRFORCE','Laughlin AFB','TX','78843');
SELECT insert_duty_station('Little Rock AFB','AIRFORCE','Little Rock AFB','AR','72099');
SELECT insert_duty_station('MacDill AFB','AIRFORCE','Tampa','FL','33621');
SELECT insert_duty_station('Maxwell AFB','AIRFORCE','Montgomery','AL','36112');
SELECT insert_duty_station('McAlester AAP','ARMY','McAlester','OK','74501');
SELECT insert_duty_station('Moody AFB','AIRFORCE','Moody AFB','GA','31699');
SELECT insert_duty_station('NAS Corpus Christi','NAVY','Corpus Christi','TX','78419');
SELECT insert_duty_station('NAS Fort Worth','NAVY','Fort Worth','TX','76127');
SELECT insert_duty_station('NAVSUP FLC Panama City','NAVY','Panama City Beach','FL','32407');
SELECT insert_duty_station('NCBC Gulfport','NAVY','Gulfport','MS','39501');
SELECT insert_duty_station('Patrick AFB','AIRFORCE','Patrick AFB','FL','32925');
SELECT insert_duty_station('Pine Bluff Arsenal','ARMY','White Hall','AR','71602');
SELECT insert_duty_station('Randolph AFB','AIRFORCE','Randolph AFB','TX','78150');
SELECT insert_duty_station('Red River Army Depot','ARMY','Texarkana','TX','75507');
SELECT insert_duty_station('Robins AFB','AIRFORCE','Warner Robins','GA','31098');
SELECT insert_duty_station('Sheppard AFB','AIRFORCE','Sheppard AFB','TX','76311');
SELECT insert_duty_station('Tinker AFB','AIRFORCE','Oklahoma City','OK','73145');
SELECT insert_duty_station('Tyndall AFB','AIRFORCE','Panama City','FL','32403');
SELECT insert_duty_station('USCG Mobile','COASTGUARD','Mobile','AL','36615');
SELECT insert_duty_station('Vance AFB','AIRFORCE','Enid','OK','73705');
SELECT insert_duty_station('Beale AFB','AIRFORCE','Beale AFB','CA','95903');
SELECT insert_duty_station('Blue Grass Army Depot','ARMY','Richmond','KY','40475');
SELECT insert_duty_station('Buckley AFB','AIRFORCE','Aurora','CO','80011');
SELECT insert_duty_station('Cannon AFB','AIRFORCE','Cannon AFB','NM','88103');
SELECT insert_duty_station('Creech AFB','AIRFORCE','Indian Springs','NV','89018');
SELECT insert_duty_station('Davis Monthan AFB','AIRFORCE','Tucson','AZ','85707');
SELECT insert_duty_station('Dugway Proving Grounds','ARMY','Dugway','UT','84022');
SELECT insert_duty_station('Edwards AFB','AIRFORCE','Edwards','CA','93524');
SELECT insert_duty_station('Ellsworth AFB','AIRFORCE','Ellsworth AFB','SD','57706');
SELECT insert_duty_station('F.E. Warren AFB','AIRFORCE','F.E. Warren AFB','WY','82005');
SELECT insert_duty_station('Fairchild AFB','AIRFORCE','Spokane','WA','99208');
SELECT insert_duty_station('Ft Carson','ARMY','Colorado Springs','CO','80913');
SELECT insert_duty_station('Ft McCoy','ARMY','Sparta','WI','54656');
SELECT insert_duty_station('Grand Forks AFB','AIRFORCE','Grand Forks AFB','ND','58205');
SELECT insert_duty_station('Hill AFB','AIRFORCE','Hill AFB','UT','84056');
SELECT insert_duty_station('Holloman AFB','AIRFORCE','Holloman AFB','NM','88330');
SELECT insert_duty_station('Kirtland AFB','AIRFORCE','Kortland AFB','NM','87117');
SELECT insert_duty_station('Luke AFB','AIRFORCE','Glendale Luke AFB','AZ','85309');
SELECT insert_duty_station('Malmstrom AFB','AIRFORCE','Malmstrom AFB','MT','59402');
SELECT insert_duty_station('McConnell AFB','AIRFORCE','McConnell AFB','KS','67221');
SELECT insert_duty_station('MCMWTC Bridgeport','MARINES','Bridgeport','CA','93517');
SELECT insert_duty_station('Minot AFB','AIRFORCE','Minot AFB','ND','58705');
SELECT insert_duty_station('Mountain Home AFB','AIRFORCE','Mountain Home AFB','ID','83648');
SELECT insert_duty_station('Nellis AFB','AIRFORCE','Nellis AFB','NV','89191');
SELECT insert_duty_station('Offutt AFB','AIRFORCE','Offutt AFB','NE','68113');
SELECT insert_duty_station('Peterson AFB','AIRFORCE','Colorado Springs','CO','80916');
SELECT insert_duty_station('Schiever AFB','AIRFORCE','Colorado Springs','CO','80912');
SELECT insert_duty_station('Tooele Army Depot','ARMY','Tooele','UT','84074');
SELECT insert_duty_station('Travis AFB','AIRFORCE','Travis AFB','CA','94535');
SELECT insert_duty_station('USAF Academy','AIRFORCE','USAF Academy','CO','80840');
SELECT insert_duty_station('USCG Humboldt Bay','COASTGUARD','Samoa','CA','95564');
SELECT insert_duty_station('Vandenberg AFB','AIRFORCE','Lompoc','CA','93437');
SELECT insert_duty_station('Whiteman AFB','AIRFORCE','Whiteman AFB','MO','65305');
SELECT insert_duty_station('White Sands Missile Range','ARMY','White Sands Missile Range','NM','88002');

DROP FUNCTION insert_duty_station;
