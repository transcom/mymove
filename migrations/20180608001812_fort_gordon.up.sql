-- We have duplicate Fort Gordon Transportation Offices in the DB
-- The next few lines delete one with the wrong address
delete from office_emails where id='b452f7b1-036f-4650-bf11-25523c622958';
delete from office_phone_lines where id='f229d042-fa19-458e-ad27-493563ab5465';
delete from office_phone_lines where id='ed5aadd1-e068-4be6-9097-2c3d75531b7b';
delete from transportation_offices where id='852d03aa-f816-4183-9b7e-a5dc2aa499e4';
delete from addresses where id='428dae6c-17aa-4403-b4ed-dc851d903561';

-- Now add Fort Gordon to the list of Duty Stations
INSERT INTO addresses VALUES ('5ac95be8-0230-47ea-90b4-b0f6f60de364',  'Fort Gordon', NULL, 'Augusta', 'GA', '30813', now(), now(), NULL, 'United States');
INSERT INTO duty_stations VALUES ('2d5ada83-e09a-47f8-8de6-83ec51694a86', 'Fort Gordon', 'ARMY', '5ac95be8-0230-47ea-90b4-b0f6f60de364',now(), now(), '19bd6cfc-35a9-4aa4-bbff-dd5efa7a9e3f');
