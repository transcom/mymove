-- Generated programmatically by load-transportation-offices.py

-- Update the TO
UPDATE transportation_offices SET name = 'CPPSO Fort Cavazos (HBAT) - USA', gbloc = 'HBAT' WHERE id = 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b';

-- Update the address
UPDATE addresses SET street_address_1 = '18010 TJ Mills Blvd', street_address_2 = 'Copeland Soldier Service Center, Room A104', city = 'Fort Cavazos', state = 'TX', postal_code = '76544' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Liberty - USA', gbloc = 'BGAC' WHERE id = 'e3c44c50-ece0-4692-8aad-8c9600000000';

-- Update the address
UPDATE addresses SET street_address_1 = 'Installation Transportation Office', street_address_2 = 'Bldg 4-2843 Normandy Drive', city = 'Fort Liberty', state = 'NC', postal_code = '28310' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e3c44c50-ece0-4692-8aad-8c9600000000');

-- Update the TO
UPDATE transportation_offices SET name = 'CPPSO Norfolk (BGNC) - USN', gbloc = 'BGNC' WHERE id = '5f741385-0a34-4d05-9068-e1e2dd8dfefc';

-- Update the address
UPDATE addresses SET street_address_1 = '7920 14th St', street_address_2 = 'Bldg SDA-336', city = 'Norfolk', state = 'VA', postal_code = '23510' WHERE id = (SELECT address_id FROM transportation_offices where id = '5f741385-0a34-4d05-9068-e1e2dd8dfefc');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - Mid Atlantic (BGAC) - USA', gbloc = 'BGAC' WHERE id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a50';

-- Update the address
UPDATE addresses SET street_address_1 = '10109 Gridley Rd ', street_address_2 = 'Bldg 314', city = 'Fort Belvoir', state = 'VA', postal_code = '22060' WHERE id = (SELECT address_id FROM transportation_offices where id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a50');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - North Central (KKFA) - USAF', gbloc = 'KKFA' WHERE id = '171b54fa-4c89-45d8-8111-a2d65818ff8c';

-- Update the address
UPDATE addresses SET street_address_1 = '121 S. Tejon St', street_address_2 = 'Suite 800', city = 'Colorado Springs', state = 'CO', postal_code = '80903' WHERE id = (SELECT address_id FROM transportation_offices where id = '171b54fa-4c89-45d8-8111-a2d65818ff8c');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - North East (AGFM) - USAF', gbloc = 'AGFM' WHERE id = '3132b512-1889-4776-a666-9c08a24afe20';

-- Update the address
UPDATE addresses SET street_address_1 = '25 Chennault St ', street_address_2 = 'Bldg 1723', city = 'Hanscom AFB', state = 'MA', postal_code = '01731' WHERE id = (SELECT address_id FROM transportation_offices where id = '3132b512-1889-4776-a666-9c08a24afe20');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - North West (JEAT) - USA', gbloc = 'JEAT' WHERE id = '5a3388e1-6d46-4639-ac8f-a8937dc26938';

-- Update the address
UPDATE addresses SET street_address_1 = '9501 Rainier Dr', street_address_2 = 'Bldg 9501', city = 'JB Lewis-McChord', state = 'WA', postal_code = '98433' WHERE id = (SELECT address_id FROM transportation_offices where id = '5a3388e1-6d46-4639-ac8f-a8937dc26938');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - South Central (HAFC) - USAF', gbloc = 'HAFC' WHERE id = 'c2c440ae-5394-4483-84fb-f872e32126bb';

-- Update the address
UPDATE addresses SET street_address_1 = '2261 Hughes Ave', street_address_2 = 'Suite 160', city = 'Lackland AFB', state = 'TX', postal_code = '78236' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c2c440ae-5394-4483-84fb-f872e32126bb');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - South East (CNNQ) - USN', gbloc = 'CNNQ' WHERE id = 'aa899628-dabb-4724-8e4a-b4579c1550e0';

-- Update the address
UPDATE addresses SET street_address_1 = 'NAVSUP FLCJ Code 4031', street_address_2 = 'PO Box 97, Bldg 110-1', city = 'Jacksonville', state = 'FL', postal_code = '32212' WHERE id = (SELECT address_id FROM transportation_offices where id = 'aa899628-dabb-4724-8e4a-b4579c1550e0');

-- Update the TO
UPDATE transportation_offices SET name = 'JPPSO - South West (LKNQ) - USN', gbloc = 'LKNQ' WHERE id = '27002d34-e9ea-4ef5-a086-f23d07c4088c';

-- Update the address
UPDATE addresses SET street_address_1 = '3376 Albacore Alley', street_address_2 = 'Bldg 3376', city = 'San Diego', state = 'CA', postal_code = '92136' WHERE id = (SELECT address_id FROM transportation_offices where id = '27002d34-e9ea-4ef5-a086-f23d07c4088c');

-- Update the TO
UPDATE transportation_offices SET name = 'Personal Property Activity HQ (PPA HQ) - USAF', gbloc = 'HAFC' WHERE id = 'ebdacf64-353a-4014-91db-0d04d88320f0';

-- Update the address
UPDATE addresses SET street_address_1 = '1960 1st Street West', street_address_2 = 'Bldg 977, Room C-101', city = 'JBSA Randolph', state = 'TX', postal_code = '78150' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ebdacf64-353a-4014-91db-0d04d88320f0');

-- Update the TO
UPDATE transportation_offices SET name = 'PPM Processing (FINCEN) - USCG', gbloc = 'USCG' WHERE id = '6598f143-9451-4635-b921-94a4f86f9ed1';

-- Update the address
UPDATE addresses SET street_address_1 = 'P.O. Box 4102', street_address_2 = 'nan', city = 'Chesapeake', state = 'VA', postal_code = '23327' WHERE id = (SELECT address_id FROM transportation_offices where id = '6598f143-9451-4635-b921-94a4f86f9ed1');

-- Update the TO
UPDATE transportation_offices SET name = 'PPM Processing (TVCB) - USMC ', gbloc = 'TVCB' WHERE id = 'aed72dfd-41cf-4111-9131-f27066aa9e88';

-- Update the address
UPDATE addresses SET street_address_1 = '814 Radford Blvd', street_address_2 = 'Suite 20262', city = 'Albany', state = 'GA', postal_code = '31704' WHERE id = (SELECT address_id FROM transportation_offices where id = 'aed72dfd-41cf-4111-9131-f27066aa9e88');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO - Outbound NAS Patuxent - USN', gbloc = 'BGNC' WHERE id = 'c12dcaae-ffcd-44aa-ae7b-2798f9e5f418';

-- Update the address
UPDATE addresses SET street_address_1 = '47253 Whalen Rd', street_address_2 = 'Bldg 588, Room 109', city = 'NAS Patuxent River', state = 'MD', postal_code = '20670' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c12dcaae-ffcd-44aa-ae7b-2798f9e5f418');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Aberdeen Proving Ground - USA', gbloc = 'BGAC' WHERE id = '6a27dfbd-2a49-485f-86dd-49475d5facef';

-- Update the address
UPDATE addresses SET street_address_1 = '4305 Susquehanna Ave', street_address_2 = 'Room 307', city = 'Aberdeen Proving Ground', state = 'MD', postal_code = '21005' WHERE id = (SELECT address_id FROM transportation_offices where id = '6a27dfbd-2a49-485f-86dd-49475d5facef');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Altus AFB - USAF', gbloc = 'HAFC' WHERE id = '3be2381f-f9ed-4902-bbdc-69c69e43eb86';

-- Update the address
UPDATE addresses SET street_address_1 = '308 N. 1st St', street_address_2 = '97 LRS/LGRDF, Bldg 52 Room 2601', city = 'Altus AFB', state = 'OK', postal_code = '73523' WHERE id = (SELECT address_id FROM transportation_offices where id = '3be2381f-f9ed-4902-bbdc-69c69e43eb86');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Aviation Training Center Elizabeth City - USCG', gbloc = 'BGNC' WHERE id = 'f7d9f4a4-c097-4c72-b0a8-41fc59e9cf44';

-- Update the address
UPDATE addresses SET street_address_1 = '1664 Weeksville Rd ', street_address_2 = 'ATTN: Transportation Office ', city = 'Elizabeth City', state = 'NC', postal_code = '27909' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f7d9f4a4-c097-4c72-b0a8-41fc59e9cf44');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Aviation Training Center Mobile - USCG', gbloc = 'HAFC' WHERE id = '27c8d71a-c3c6-4a2d-ac1b-874cbf6a5f85';

-- Update the address
UPDATE addresses SET street_address_1 = '8501 Tanner Williams Rd', street_address_2 = 'nan', city = 'Mobile', state = 'AL', postal_code = '36608' WHERE id = (SELECT address_id FROM transportation_offices where id = '27c8d71a-c3c6-4a2d-ac1b-874cbf6a5f85');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Barksdale AFB - USAF', gbloc = 'HAFC' WHERE id = '86ef23b7-c4b7-4cd0-af7a-29570ab8fa84';

-- Update the address
UPDATE addresses SET street_address_1 = '801 Kenny Ave', street_address_2 = 'Suite 3400', city = 'Barksdale AFB', state = 'LA', postal_code = '71110' WHERE id = (SELECT address_id FROM transportation_offices where id = '86ef23b7-c4b7-4cd0-af7a-29570ab8fa84');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Alameda - USCG', gbloc = 'LHNQ' WHERE id = '3fc4b408-1197-430a-a96a-24a5a1685b45';

-- Update the address
UPDATE addresses SET street_address_1 = 'Coast Guard Island', street_address_2 = 'Bldg 3', city = 'Alameda', state = 'CA', postal_code = '94501' WHERE id = (SELECT address_id FROM transportation_offices where id = '3fc4b408-1197-430a-a96a-24a5a1685b45');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Boston - USCG', gbloc = 'AGFM' WHERE id = 'cf1addea-a4f9-4173-8506-2bb82a064cb7';

-- Update the address
UPDATE addresses SET street_address_1 = 'ATTN: Transportation Office ', street_address_2 = '427 Commerical Street ', city = 'Boston', state = 'MA', postal_code = '02109' WHERE id = (SELECT address_id FROM transportation_offices where id = 'cf1addea-a4f9-4173-8506-2bb82a064cb7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Cape Cod - USCG', gbloc = 'AGFM' WHERE id = 'ca40217d-c4e0-4931-b181-e8b99c4a2a75';

-- Update the address
UPDATE addresses SET street_address_1 = '3159 Herbert Rd', street_address_2 = 'nan', city = 'Buzzards Bay', state = 'MA', postal_code = '02542' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ca40217d-c4e0-4931-b181-e8b99c4a2a75');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Charleston - USCG', gbloc = 'AGFM' WHERE id = 'b429ec35-eb8b-4c35-979b-e96cb9a9cbf7';

-- Update the address
UPDATE addresses SET street_address_1 = '1050 Register St', street_address_2 = 'nan', city = 'North Charleston', state = 'SC', postal_code = '29405' WHERE id = (SELECT address_id FROM transportation_offices where id = 'b429ec35-eb8b-4c35-979b-e96cb9a9cbf7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Los Angeles/Long Beach - USCG', gbloc = 'LKNQ' WHERE id = 'f3ce3548-fb03-4e2e-b986-386a1adf2400';

-- Update the address
UPDATE addresses SET street_address_1 = '1001 South Seaside Ave', street_address_2 = 'nan', city = 'San Pedro', state = 'CA', postal_code = '90731' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f3ce3548-fb03-4e2e-b986-386a1adf2400');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Miami - USCG', gbloc = 'CLPK' WHERE id = '1b3e7496-efa7-48aa-ba22-b630d6fea98b';

-- Update the address
UPDATE addresses SET street_address_1 = '15610 SW 117th Ave ', street_address_2 = 'nan', city = 'Miami', state = 'FL', postal_code = '33177' WHERE id = (SELECT address_id FROM transportation_offices where id = '1b3e7496-efa7-48aa-ba22-b630d6fea98b');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base National Capital Region - USCG', gbloc = 'BGAC' WHERE id = '12377a7a-7cd0-4c75-acbb-1a19242909f0';

-- Update the address
UPDATE addresses SET street_address_1 = '2703 Martin Luther King Ave SE', street_address_2 = 'Stop 7118', city = 'Washington', state = 'DC', postal_code = '20593' WHERE id = (SELECT address_id FROM transportation_offices where id = '12377a7a-7cd0-4c75-acbb-1a19242909f0');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Portsmouth - USCG', gbloc = 'BGNC' WHERE id = '3021df82-cb36-4286-bfb6-94051d04b59b';

-- Update the address
UPDATE addresses SET street_address_1 = '4000 Coast Guard Blvd', street_address_2 = 'nan', city = 'Portsmouth', state = 'VA', postal_code = '23707' WHERE id = (SELECT address_id FROM transportation_offices where id = '3021df82-cb36-4286-bfb6-94051d04b59b');

-- Insert the address
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('6ab7be64-1499-486e-b4ae-2c1a8f3f069d', '1001 S. Seaside Ave', 'Bldg 25', 'San Pedro', 'CA', '90731', now(), now());

-- Insert the TO
INSERT INTO transportation_offices (id, address_id, gbloc, name, created_at, updated_at, latitude, longitude) VALUES ('f0ccc12b-ec87-4767-8053-bf35409a3910', '6ab7be64-1499-486e-b4ae-2c1a8f3f069d', 'KKFA', 'PPPO Base San Pedro - USCG', now(), now(), 0, 0);

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base Seattle - USCG', gbloc = 'JEAT' WHERE id = '183969ce-8abd-4136-b193-2041a8c4f1be';

-- Update the address
UPDATE addresses SET street_address_1 = 'ATTN: Transportation Office ', street_address_2 = '1519 Alaskan Way South', city = 'Seattle', state = 'WA', postal_code = '98134' WHERE id = (SELECT address_id FROM transportation_offices where id = '183969ce-8abd-4136-b193-2041a8c4f1be');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Base St Louis - USCG', gbloc = 'AGFM' WHERE id = '254765b6-1781-4441-b7f1-27c42beaa583';

-- Update the address
UPDATE addresses SET street_address_1 = '1222 Spruce St, Suite 2.107 (Ray Federal Bldg)', street_address_2 = 'Room 2.107', city = 'St Louis', state = 'MO', postal_code = '63103' WHERE id = (SELECT address_id FROM transportation_offices where id = '254765b6-1781-4441-b7f1-27c42beaa583');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Beale AFB - USAF', gbloc = 'KKFA' WHERE id = '9685f46b-b288-4bb3-b2f8-b70a23b91943';

-- Update the address
UPDATE addresses SET street_address_1 = '17852 16th St', street_address_2 = 'Bldg 25216', city = 'Beale AFB', state = 'CA', postal_code = '95903' WHERE id = (SELECT address_id FROM transportation_offices where id = '9685f46b-b288-4bb3-b2f8-b70a23b91943');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Blue Grass Army Depot - USA', gbloc = 'KKFA' WHERE id = '175a4999-ccf7-4077-b25f-955e19fd5873';

-- Update the address
UPDATE addresses SET street_address_1 = '431 Battlefield Memorial Hwy', street_address_2 = 'Bldg 15', city = 'Richmond', state = 'KY', postal_code = '40475' WHERE id = (SELECT address_id FROM transportation_offices where id = '175a4999-ccf7-4077-b25f-955e19fd5873');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Buckley SFB - USAF', gbloc = 'KKFA' WHERE id = '9ff34ab7-2e87-4515-be3a-f35c90e65e6a';

-- Update the address
UPDATE addresses SET street_address_1 = '18401 E. A Basin Ave', street_address_2 = 'Bldg 606', city = 'Buckley SFB', state = 'CO', postal_code = '80011' WHERE id = (SELECT address_id FROM transportation_offices where id = '9ff34ab7-2e87-4515-be3a-f35c90e65e6a');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Cannon AFB - USAF', gbloc = 'KKFA' WHERE id = '80796bc4-e494-4b19-bb16-cdcdba187829';

-- Update the address
UPDATE addresses SET street_address_1 = '110 E. Cochran Ave', street_address_2 = 'Bldg 600 Room 1153', city = 'Cannon AFB', state = 'NM', postal_code = '88103' WHERE id = (SELECT address_id FROM transportation_offices where id = '80796bc4-e494-4b19-bb16-cdcdba187829');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Carlisle Barracks - USA', gbloc = 'BGAC' WHERE id = 'e37860be-c642-4037-af9a-8a1be690d8d7';

-- Update the address
UPDATE addresses SET street_address_1 = '46 Ashburn Dr Room 227-232', street_address_2 = '2nd Floor', city = 'Carlisle', state = 'PA', postal_code = '17013' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e37860be-c642-4037-af9a-8a1be690d8d7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Columbus AFB - USAF', gbloc = 'HAFC' WHERE id = '1ef624a7-d612-4457-b310-d9e4c3fd6c85';

-- Update the address
UPDATE addresses SET street_address_1 = '495 Harpe Blvd', street_address_2 = 'Bldg 730, Suite 120', city = 'Columbus AFB', state = 'MS', postal_code = '39710' WHERE id = (SELECT address_id FROM transportation_offices where id = '1ef624a7-d612-4457-b310-d9e4c3fd6c85');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Corpus Christi - USN', gbloc = 'HAFC' WHERE id = '5853ef41-0169-4050-9483-2d2bb382bd4e';

-- Update the address
UPDATE addresses SET street_address_1 = '9035 Ocean Dr', street_address_2 = '3rd Floor', city = 'Corpus Christi', state = 'TX', postal_code = '78419' WHERE id = (SELECT address_id FROM transportation_offices where id = '5853ef41-0169-4050-9483-2d2bb382bd4e');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Creech AFB - USAF', gbloc = 'KKFA' WHERE id = '5de30a80-a8e5-458c-9b54-edfae7b8cdb9';

-- Update the address
UPDATE addresses SET street_address_1 = '1012 Perimeter Rd', street_address_2 = 'Bldg 1012', city = 'Creech AFB (Indian Springs)', state = 'NV', postal_code = '89018' WHERE id = (SELECT address_id FROM transportation_offices where id = '5de30a80-a8e5-458c-9b54-edfae7b8cdb9');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Davis Monthan AFB - USAF', gbloc = 'KKFA' WHERE id = '54156892-dff1-4657-8998-39ff4e3a259e';

-- Update the address
UPDATE addresses SET street_address_1 = '3515 S Fifth St', street_address_2 = 'nan', city = 'Davis-Monthan AFB', state = 'AZ', postal_code = '85707' WHERE id = (SELECT address_id FROM transportation_offices where id = '54156892-dff1-4657-8998-39ff4e3a259e');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Detroit Arsenal - USA', gbloc = 'BGAC' WHERE id = '07c5e487-ee37-4e43-b2db-2b4bd79ea7e5';

-- Update the address
UPDATE addresses SET street_address_1 = '6501 E. 11 Mile Rd', street_address_2 = 'Bldg 232', city = 'Detroit Arsenal', state = 'MI', postal_code = '48397' WHERE id = (SELECT address_id FROM transportation_offices where id = '07c5e487-ee37-4e43-b2db-2b4bd79ea7e5');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO Albany - USMC ', gbloc = 'CNNQ' WHERE id = '65bc635c-c097-428b-a4e5-8b752510f22e';

-- Update the address
UPDATE addresses SET street_address_1 = '814 Radford Blvd', street_address_2 = 'Bldg 3500, Wing 500, Rooms 501 and 503', city = 'Albany', state = 'GA', postal_code = '31704' WHERE id = (SELECT address_id FROM transportation_offices where id = '65bc635c-c097-428b-a4e5-8b752510f22e');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO Camp Lejeune - USMC ', gbloc = 'CNNQ' WHERE id = '22894aa1-1c29-49d8-bd1b-2ce64448cc8d';

-- Update the address
UPDATE addresses SET street_address_1 = 'Ash St', street_address_2 = 'Bldg 1011', city = 'Camp Lejeune', state = 'NC', postal_code = '28547' WHERE id = (SELECT address_id FROM transportation_offices where id = '22894aa1-1c29-49d8-bd1b-2ce64448cc8d');

-- Insert the address
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('ed59a5f8-4ecd-40fb-9e68-af76e9fd5362', 'Vandergrift Blvd', 'Bldg 2263', 'Camp Pendleton', 'CA', '92055', now(), now());

-- Insert the TO
INSERT INTO transportation_offices (id, address_id, gbloc, name, created_at, updated_at, latitude, longitude) VALUES ('0435ccdb-5c03-496c-b09f-b4c20e456659', 'ed59a5f8-4ecd-40fb-9e68-af76e9fd5362', 'USMC', 'PPPO DMO Camp Pendleton - USMC', now(), now(), 0, 0);

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO Camp Pendleton - USMC', gbloc = 'LKNQ' WHERE id = 'f50eb7f5-960a-46e8-aa64-6025b44132ab';

-- Update the address
UPDATE addresses SET street_address_1 = 'Vandergrift Blvd', street_address_2 = 'Bldg 2263', city = 'Camp Pendleton', state = 'CA', postal_code = '92055' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f50eb7f5-960a-46e8-aa64-6025b44132ab');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO JB Myer-Henderson Hall - USMC ', gbloc = 'BGAC' WHERE id = '20e19766-555d-486d-96a0-995d4d2cdacf';

-- Update the address
UPDATE addresses SET street_address_1 = '1555 Southgate Rd', street_address_2 = 'Henderson Hall, Bldg 29, Room 302', city = 'Arlington', state = 'VA', postal_code = '22214' WHERE id = (SELECT address_id FROM transportation_offices where id = '20e19766-555d-486d-96a0-995d4d2cdacf');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO MCAGCC 29 Palms - USMC', gbloc = 'LKNQ' WHERE id = 'bd733387-6b6c-42ba-b2c3-76c20cc65006';

-- Update the address
UPDATE addresses SET street_address_1 = 'DMO Personal Property Office', street_address_2 = 'Bldg 1102, Door 22, Marine Corps Base', city = '29 Palms', state = 'CA', postal_code = '92278' WHERE id = (SELECT address_id FROM transportation_offices where id = 'bd733387-6b6c-42ba-b2c3-76c20cc65006');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO MCAS Beaufort - USMC', gbloc = 'CNNQ' WHERE id = 'e83c17ae-600c-43aa-ba0b-321936038f36';

-- Update the address
UPDATE addresses SET street_address_1 = 'Drayton St', street_address_2 = 'Bldg 612', city = 'MCAS Beaufort', state = 'SC', postal_code = '29904' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e83c17ae-600c-43aa-ba0b-321936038f36');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO MCAS Cherry Point - USMC', gbloc = 'CNNQ' WHERE id = '1c17dc49-f411-4815-9b96-71b26a960f7b';

-- Update the address
UPDATE addresses SET street_address_1 = 'E Street', street_address_2 = 'Bldg 298', city = 'Cherry Point', state = 'NC', postal_code = '28533' WHERE id = (SELECT address_id FROM transportation_offices where id = '1c17dc49-f411-4815-9b96-71b26a960f7b');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO MCAS Miramar - USMC', gbloc = 'LKNQ' WHERE id = '42f7ef32-7d1f-4c03-b0df-6af9832615fc';

-- Update the address
UPDATE addresses SET street_address_1 = '2258 Delta Rd', street_address_2 = 'nan', city = 'Miramar', state = 'CA', postal_code = '92145' WHERE id = (SELECT address_id FROM transportation_offices where id = '42f7ef32-7d1f-4c03-b0df-6af9832615fc');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO MCAS Yuma - USMC', gbloc = 'LKNQ' WHERE id = '6ac7e595-1e0c-44cb-a9a4-cd7205868ed4';

-- Update the address
UPDATE addresses SET street_address_1 = 'Spears St', street_address_2 = 'Bldg 328, 2nd Floor', city = 'Yuma', state = 'AZ', postal_code = '85369' WHERE id = (SELECT address_id FROM transportation_offices where id = '6ac7e595-1e0c-44cb-a9a4-cd7205868ed4');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO MCB Quantico - USMC ', gbloc = 'BGAC' WHERE id = '2ffbe627-9918-4f52-a440-4be87f5fca73';

-- Update the address
UPDATE addresses SET street_address_1 = '2009 Zeilin Rd', street_address_2 = 'PPPO DMO Quantico, Bldg 2009', city = 'MCB Quantico', state = 'VA', postal_code = '22134' WHERE id = (SELECT address_id FROM transportation_offices where id = '2ffbe627-9918-4f52-a440-4be87f5fca73');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO DMO Mountain Warfare Training Center Bridgeport - USMC ', gbloc = 'KKFA' WHERE id = 'fab58a38-ee1f-4adf-929a-2dd246fc5e67';

-- Update the address
UPDATE addresses SET street_address_1 = 'HC 83', street_address_2 = 'Bldg 7047', city = 'Bridgeport', state = 'CA', postal_code = '93517' WHERE id = (SELECT address_id FROM transportation_offices where id = 'fab58a38-ee1f-4adf-929a-2dd246fc5e67');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Dover AFB - USAF', gbloc = 'AGFM' WHERE id = '3a43dc63-be80-40ff-8410-839e6658e35c';

-- Update the address
UPDATE addresses SET street_address_1 = '550 Atlantic St', street_address_2 = 'nan', city = 'Dover AFB', state = 'DE', postal_code = '19902' WHERE id = (SELECT address_id FROM transportation_offices where id = '3a43dc63-be80-40ff-8410-839e6658e35c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Dyess AFB - USAF', gbloc = 'HAFC' WHERE id = '4e28dc64-3368-4ab4-b98b-4c8d60aacc2d';

-- Update the address
UPDATE addresses SET street_address_1 = '417 3rd St', street_address_2 = 'Bldg 7402', city = 'Dyess AFB', state = 'TX', postal_code = '79607' WHERE id = (SELECT address_id FROM transportation_offices where id = '4e28dc64-3368-4ab4-b98b-4c8d60aacc2d');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Edwards AFB - USAF', gbloc = 'KKFA' WHERE id = '311b5292-6a8c-4ed4-a7e1-374734118737';

-- Update the address
UPDATE addresses SET street_address_1 = '5 Steller Ave', street_address_2 = 'Bldg 3000', city = 'Edwards AFB', state = 'CA', postal_code = '93524' WHERE id = (SELECT address_id FROM transportation_offices where id = '311b5292-6a8c-4ed4-a7e1-374734118737');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Eglin AFB - USAF', gbloc = 'HAFC' WHERE id = '4b8584e4-0db6-4554-9a8f-3062456884e8';

-- Update the address
UPDATE addresses SET street_address_1 = '310 Van Matre Blvd', street_address_2 = 'Bldg 210, Room 169', city = 'Eglin AFB', state = 'FL', postal_code = '32542' WHERE id = (SELECT address_id FROM transportation_offices where id = '4b8584e4-0db6-4554-9a8f-3062456884e8');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Ellington Field ANGB - USANG', gbloc = 'HAFC' WHERE id = '4716d536-7879-4eca-a5a8-cdff3a451367';

-- Update the address
UPDATE addresses SET street_address_1 = '14657 Sneider St', street_address_2 = 'Bldg 1057, Room 151', city = 'Houston', state = 'TX', postal_code = '77034' WHERE id = (SELECT address_id FROM transportation_offices where id = '4716d536-7879-4eca-a5a8-cdff3a451367');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Ellsworth AFB - USAF', gbloc = 'KKFA' WHERE id = '07776e49-5100-4094-a7b9-5d9de4fa18d6';

-- Update the address
UPDATE addresses SET street_address_1 = '1924 Lemay Blvd', street_address_2 = 'nan', city = 'Ellsworth AFB', state = 'SD', postal_code = '57706' WHERE id = (SELECT address_id FROM transportation_offices where id = '07776e49-5100-4094-a7b9-5d9de4fa18d6');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fairchild AFB - USAF', gbloc = 'KKFA' WHERE id = '972c238d-89a0-4b50-a0cf-79995c3ed1e7';

-- Update the address
UPDATE addresses SET street_address_1 = '220 W Bong St', street_address_2 = 'Bldg 2045, Room 127', city = 'Fairchild AFB', state = 'WA', postal_code = '99011' WHERE id = (SELECT address_id FROM transportation_offices where id = '972c238d-89a0-4b50-a0cf-79995c3ed1e7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO FE Warren AFB - USAF', gbloc = 'KKFA' WHERE id = '485ab35a-79e0-4db4-9f13-09f57532deee';

-- Update the address
UPDATE addresses SET street_address_1 = '7100 Saber Rd', street_address_2 = 'Bldg 1284', city = 'FE Warren AFB', state = 'WY', postal_code = '82001' WHERE id = (SELECT address_id FROM transportation_offices where id = '485ab35a-79e0-4db4-9f13-09f57532deee');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO FLCJ DET New Orleans - USN', gbloc = 'CNNQ' WHERE id = '358ab6bb-9d2c-4f07-9be8-b69e1a92b4f8';

-- Update the address
UPDATE addresses SET street_address_1 = '400 Russell Ave', street_address_2 = 'Bldg 31', city = 'New Orleans', state = 'LA', postal_code = '70143' WHERE id = (SELECT address_id FROM transportation_offices where id = '358ab6bb-9d2c-4f07-9be8-b69e1a92b4f8');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO FLCJ Meridan Det - USN', gbloc = 'CNNQ' WHERE id = 'f8c700ae-6633-4092-95a5-dddbf10da356';

-- Update the address
UPDATE addresses SET street_address_1 = '224 Allen Rd', street_address_2 = 'NAS Meridian Bldg', city = 'Meridian', state = 'MS', postal_code = '39302' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f8c700ae-6633-4092-95a5-dddbf10da356');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Bliss - USA', gbloc = 'HBAT' WHERE id = '50579f6f-b23a-4d6f-a4c6-62961f09f7a7';

-- Update the address
UPDATE addresses SET street_address_1 = '504A Holbrook Rd', street_address_2 = '1st Floor, Room 106', city = 'Fort Bliss', state = 'TX', postal_code = '79916' WHERE id = (SELECT address_id FROM transportation_offices where id = '50579f6f-b23a-4d6f-a4c6-62961f09f7a7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Campbell - USA', gbloc = 'BGAC' WHERE id = '6ea40690-4f05-4c06-8775-82ed0f160d47';

-- Update the address
UPDATE addresses SET street_address_1 = '7162 Hedgerow Rd', street_address_2 = 'nan', city = 'Fort Campbell', state = 'KY', postal_code = '42223' WHERE id = (SELECT address_id FROM transportation_offices where id = '6ea40690-4f05-4c06-8775-82ed0f160d47');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Carson - USAF', gbloc = 'KKFA' WHERE id = 'c696167c-9320-4c73-b0ce-80ec0920d9fd';

-- Update the address
UPDATE addresses SET street_address_1 = '6351 Wetzel Ave', street_address_2 = 'Bldg 1525', city = 'Fort Carson', state = 'CO', postal_code = '80913' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c696167c-9320-4c73-b0ce-80ec0920d9fd');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Detrick - USA', gbloc = 'BGAC' WHERE id = '1fd1720c-4dfb-40e4-b661-b4ac994deaae';

-- Update the address
UPDATE addresses SET street_address_1 = '9200 Amber Dr', street_address_2 = 'Room 218', city = 'Fort Detrick', state = 'MD', postal_code = '21702' WHERE id = (SELECT address_id FROM transportation_offices where id = '1fd1720c-4dfb-40e4-b661-b4ac994deaae');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Drum - USA', gbloc = 'BGAC' WHERE id = '45aff62e-3f27-478e-a7ab-db8fecb8ac2e';

-- Update the address
UPDATE addresses SET street_address_1 = '10720 Belvedere Blvd', street_address_2 = 'Room A2-42', city = 'Fort Drum', state = 'NY', postal_code = '13602' WHERE id = (SELECT address_id FROM transportation_offices where id = '45aff62e-3f27-478e-a7ab-db8fecb8ac2e');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Eisenhower - USA', gbloc = 'CNNQ' WHERE id = '19bd6cfc-35a9-4aa4-bbff-dd5efa7a9e3f';

-- Update the address
UPDATE addresses SET street_address_1 = '307 Chamberlain Ave', street_address_2 = 'Bldg 33720, Room 111A', city = 'Fort Eisenhower', state = 'GA', postal_code = '30905' WHERE id = (SELECT address_id FROM transportation_offices where id = '19bd6cfc-35a9-4aa4-bbff-dd5efa7a9e3f');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Gregg-Adams - USA', gbloc = 'BGAC' WHERE id = '4cc26e01-f0ea-4048-8081-1d179426a6d9';

-- Update the address
UPDATE addresses SET street_address_1 = '1401 B Ave', street_address_2 = 'Bldg 3400, Room 119', city = 'Fort Gregg-Adams', state = 'VA', postal_code = '23801' WHERE id = (SELECT address_id FROM transportation_offices where id = '4cc26e01-f0ea-4048-8081-1d179426a6d9');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Hamilton - USA', gbloc = 'BGAC' WHERE id = '9324bf4b-d84f-4a28-994f-32cdda5580d1';

-- Update the address
UPDATE addresses SET street_address_1 = '114 White Ave', street_address_2 = 'Room 114', city = 'Fort Hamilton', state = 'NY', postal_code = '11252' WHERE id = (SELECT address_id FROM transportation_offices where id = '9324bf4b-d84f-4a28-994f-32cdda5580d1');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Huachuca - USA', gbloc = 'LKNQ' WHERE id = '5c72fddb-49d3-4853-99a4-ca45f8ba07a5';

-- Update the address
UPDATE addresses SET street_address_1 = '2317 Smith Ave', street_address_2 = 'Bldg 52065', city = 'Fort Huachuca', state = 'AZ', postal_code = '85613' WHERE id = (SELECT address_id FROM transportation_offices where id = '5c72fddb-49d3-4853-99a4-ca45f8ba07a5');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Irwin - USA', gbloc = 'LKNQ' WHERE id = 'd00e3ee8-baba-4991-8f3b-86c2e370d1be';

-- Update the address
UPDATE addresses SET street_address_1 = 'Langford Lake Rd', street_address_2 = 'Bldg 105', city = 'Fort Irwin', state = 'CA', postal_code = '92310' WHERE id = (SELECT address_id FROM transportation_offices where id = 'd00e3ee8-baba-4991-8f3b-86c2e370d1be');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Jackson - USA', gbloc = 'CNNQ' WHERE id = 'a1baca88-3b40-4dff-8c87-aa38b3d5acf2';

-- Update the address
UPDATE addresses SET street_address_1 = '5450 Strom Thurmond Blvd', street_address_2 = 'Room 102', city = 'Fort Jackson', state = 'SC', postal_code = '29207' WHERE id = (SELECT address_id FROM transportation_offices where id = 'a1baca88-3b40-4dff-8c87-aa38b3d5acf2');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Johnson - USA', gbloc = 'CNNQ' WHERE id = '31fb763a-5957-4ba5-b82b-956599970b0f';

-- Update the address
UPDATE addresses SET street_address_1 = '7585 Virginia Ave', street_address_2 = 'Bldg 4374, 2nd Floor', city = 'Fort Johnson South', state = 'LA', postal_code = '71459' WHERE id = (SELECT address_id FROM transportation_offices where id = '31fb763a-5957-4ba5-b82b-956599970b0f');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Knox - USA', gbloc = 'BGAC' WHERE id = '0357f830-2f32-41f3-9ca2-268dd70df5cb';

-- Update the address
UPDATE addresses SET street_address_1 = 'LRC 25 W. Chaffee Ave', street_address_2 = 'Bldg 1384, 2nd Floor', city = 'Fort Knox', state = 'KY', postal_code = '40121' WHERE id = (SELECT address_id FROM transportation_offices where id = '0357f830-2f32-41f3-9ca2-268dd70df5cb');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Leavenworth - USA', gbloc = 'JEAT' WHERE id = 'b2f76d56-6996-41a3-aef7-483524a643d1';

-- Update the address
UPDATE addresses SET street_address_1 = '549 Kearny Ave', street_address_2 = 'Bldg 268', city = 'Fort Leavenworth', state = 'KS', postal_code = '66027' WHERE id = (SELECT address_id FROM transportation_offices where id = 'b2f76d56-6996-41a3-aef7-483524a643d1');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Leonard Wood - USA', gbloc = 'JEAT' WHERE id = 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7';

-- Update the address
UPDATE addresses SET street_address_1 = '13486 Replacement Ave', street_address_2 = 'Suite 1219/1220', city = 'Fort Leonard Wood', state = 'MO', postal_code = '65473' WHERE id = (SELECT address_id FROM transportation_offices where id = 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort McCoy - USA', gbloc = 'KKFA' WHERE id = 'dc7c6746-50b5-418a-a925-66dfa19481df';

-- Update the address
UPDATE addresses SET street_address_1 = '200 E. G St', street_address_2 = 'nan', city = 'Fort McCoy', state = 'WI', postal_code = '54656' WHERE id = (SELECT address_id FROM transportation_offices where id = 'dc7c6746-50b5-418a-a925-66dfa19481df');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Meade - USA', gbloc = 'BGAC' WHERE id = '24e187a7-ae1e-4a78-acb8-92b7e2f75950';

-- Update the address
UPDATE addresses SET street_address_1 = '4550 Parade Field Lane', street_address_2 = 'Room 131', city = 'Fort Meade', state = 'MD', postal_code = '20755' WHERE id = (SELECT address_id FROM transportation_offices where id = '24e187a7-ae1e-4a78-acb8-92b7e2f75950');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Moore - USA', gbloc = 'CNNQ' WHERE id = '5a9ed15c-ed78-47e5-8afd-7583f3cc660d';

-- Update the address
UPDATE addresses SET street_address_1 = '6650 Meloy Hall', street_address_2 = 'Bldg 6, Room 105', city = 'Fort Moore', state = 'GA', postal_code = '31905' WHERE id = (SELECT address_id FROM transportation_offices where id = '5a9ed15c-ed78-47e5-8afd-7583f3cc660d');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Novosel  - USA', gbloc = 'CNNQ' WHERE id = '2af6d201-0e8a-4837-afed-edb57ea92c4d';

-- Update the address
UPDATE addresses SET street_address_1 = '453 Novosel St', street_address_2 = 'Bldg 5700, Room 270', city = 'Fort Novosel', state = 'AL', postal_code = '36362' WHERE id = (SELECT address_id FROM transportation_offices where id = '2af6d201-0e8a-4837-afed-edb57ea92c4d');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Riley - USA', gbloc = 'JEAT' WHERE id = '86b21668-b390-495a-b4b4-ccd88e1f3ed8';

-- Update the address
UPDATE addresses SET street_address_1 = 'Custer Ave', street_address_2 = 'Bldg. 210, Room 004', city = 'Fort Riley', state = 'KS', postal_code = '66442' WHERE id = (SELECT address_id FROM transportation_offices where id = '86b21668-b390-495a-b4b4-ccd88e1f3ed8');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Sill - USA', gbloc = 'JEAT' WHERE id = '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9';

-- Update the address
UPDATE addresses SET street_address_1 = '4700 Mow Way Rd', street_address_2 = 'Room 110', city = 'Fort Sill', state = 'OK', postal_code = '73503' WHERE id = (SELECT address_id FROM transportation_offices where id = '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Fort Stewart  - USA', gbloc = 'CNNQ' WHERE id = '95b6fda3-3ce2-4fda-87df-4aefaca718c5';

-- Update the address
UPDATE addresses SET street_address_1 = '55 Pony Soldier Ave', street_address_2 = 'Bldg 253, Suite 2003A', city = 'Fort Stewart', state = 'GA', postal_code = '31314' WHERE id = (SELECT address_id FROM transportation_offices where id = '95b6fda3-3ce2-4fda-87df-4aefaca718c5');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Goodfellow AFB - USAF', gbloc = 'HAFC' WHERE id = '2c6f2ac9-210c-4053-9770-1ed6225db248';

-- Update the address
UPDATE addresses SET street_address_1 = '317 Lancaster Ave', street_address_2 = 'Bldg 423', city = 'Goodfellow AFB', state = 'TX', postal_code = '76908' WHERE id = (SELECT address_id FROM transportation_offices where id = '2c6f2ac9-210c-4053-9770-1ed6225db248');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Grand Forks AFB - USAF', gbloc = 'KKFA' WHERE id = '391ffdd2-47a5-4f6b-bda2-48babe471274';

-- Update the address
UPDATE addresses SET street_address_1 = '226 Steen Blvd', street_address_2 = 'nan', city = 'Grand Forks AFB', state = 'ND', postal_code = '58205' WHERE id = (SELECT address_id FROM transportation_offices where id = '391ffdd2-47a5-4f6b-bda2-48babe471274');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Hanscom AFB - USAF', gbloc = 'AGFM' WHERE id = '645cc264-a0ff-4ea1-9aff-84644d7ade3c';

-- Update the address
UPDATE addresses SET street_address_1 = '25 Chennault St', street_address_2 = 'Bldg 1723', city = 'Hanscom AFB', state = 'MA', postal_code = '01731' WHERE id = (SELECT address_id FROM transportation_offices where id = '645cc264-a0ff-4ea1-9aff-84644d7ade3c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Hill AFB - USAF', gbloc = 'KKFA' WHERE id = 'e05e4d21-7655-41f7-93f1-c2aac39be98c';

-- Update the address
UPDATE addresses SET street_address_1 = '7437 6th St', street_address_2 = 'Bldg 430 North', city = 'Hill AFB', state = 'UT', postal_code = '84056' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e05e4d21-7655-41f7-93f1-c2aac39be98c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Holloman AFB - USAF', gbloc = 'KKFA' WHERE id = '0bdeddf8-b539-4bc8-a86d-aa6c0c1039da';

-- Update the address
UPDATE addresses SET street_address_1 = '681 2nd St', street_address_2 = 'Bldg 222 Room 223', city = 'Holloman Air Force Base', state = 'NM', postal_code = '88330' WHERE id = (SELECT address_id FROM transportation_offices where id = '0bdeddf8-b539-4bc8-a86d-aa6c0c1039da');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Hunter Army Airfield - USA', gbloc = 'CNNQ' WHERE id = '425075d6-655e-46dc-9d0f-2dad5f0bf916';

-- Update the address
UPDATE addresses SET street_address_1 = '171 Haley Ave', street_address_2 = 'Bldg 1286, Suite 229', city = 'Hunter Army Airfield', state = 'GA', postal_code = '31406' WHERE id = (SELECT address_id FROM transportation_offices where id = '425075d6-655e-46dc-9d0f-2dad5f0bf916');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Hurlburt AFB - USAF', gbloc = 'HAFC' WHERE id = '0f4ebcf9-92c8-49bb-96f5-677e442badaa';

-- Update the address
UPDATE addresses SET street_address_1 = '212 Lukasik Ave', street_address_2 = 'Suite 157', city = 'Hurlburt Field', state = 'FL', postal_code = '32544' WHERE id = (SELECT address_id FROM transportation_offices where id = '0f4ebcf9-92c8-49bb-96f5-677e442badaa');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Jacksonville - USN', gbloc = 'CNNQ' WHERE id = 'f5ab88fe-47f8-4b58-99af-41067d6cb60d';

-- Update the address
UPDATE addresses SET street_address_1 = 'FLC Jacksonville', street_address_2 = 'Bldg 110, Room 103', city = 'Jacksonville', state = 'FL', postal_code = '32212' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f5ab88fe-47f8-4b58-99af-41067d6cb60d');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Anacostia - Bolling - USAF', gbloc = 'BGAC' WHERE id = 'd06fce02-5133-4c47-a0a4-83559826b9f5';

-- Update the address
UPDATE addresses SET street_address_1 = '229 Brookley Ave', street_address_2 = 'Bldg 520, Room 1', city = 'Washington DC', state = 'DC', postal_code = '20032' WHERE id = (SELECT address_id FROM transportation_offices where id = 'd06fce02-5133-4c47-a0a4-83559826b9f5');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Andrews AFB - USAF', gbloc = 'BGAC' WHERE id = 'd9807312-f4d0-4186-ab23-c770974ea5a7';

-- Update the address
UPDATE addresses SET street_address_1 = '1500 West Perimeter Rd', street_address_2 = 'Suite 2700', city = 'JB Andrews', state = 'MD', postal_code = '20762' WHERE id = (SELECT address_id FROM transportation_offices where id = 'd9807312-f4d0-4186-ab23-c770974ea5a7');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Charleston (Charleston) - USAF', gbloc = 'AGFM' WHERE id = 'ae98567c-3943-4ab8-92b4-771275d9b918';

-- Update the address
UPDATE addresses SET street_address_1 = '1000 Quality Cir', street_address_2 = 'Bldg 36', city = 'Goose Creek', state = 'SC', postal_code = '29445' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ae98567c-3943-4ab8-92b4-771275d9b918');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Charleston (Naval Weapon Station) - USAF', gbloc = 'AGFM' WHERE id = '67281ea0-222a-41ea-9ec2-dc274a2513ea';

-- Update the address
UPDATE addresses SET street_address_1 = '1000 Pomflant Access Rd', street_address_2 = 'Bldg 302', city = 'Goose Creek', state = 'SC', postal_code = '29445' WHERE id = (SELECT address_id FROM transportation_offices where id = '67281ea0-222a-41ea-9ec2-dc274a2513ea');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Langley-Eustis (Fort Eustis) - USAF', gbloc = 'AGFM' WHERE id = 'e3f5e683-889f-437f-b13d-3ccd7ad0d453';

-- Update the address
UPDATE addresses SET street_address_1 = '650 Monroe Ave', street_address_2 = 'nan', city = 'Fort Eustis', state = 'VA', postal_code = '23604' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e3f5e683-889f-437f-b13d-3ccd7ad0d453');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Langley-Eustis (Langley) - USAF', gbloc = 'AGFM' WHERE id = '525af4ce-3656-4ec6-a23d-2df6b4f39b64';

-- Update the address
UPDATE addresses SET street_address_1 = '45 Nealy Ave', street_address_2 = 'Suite 209', city = 'Hampton', state = 'VA', postal_code = '23665' WHERE id = (SELECT address_id FROM transportation_offices where id = '525af4ce-3656-4ec6-a23d-2df6b4f39b64');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB Lewis-McChord (Fort Lewis) - USA', gbloc = 'JEAT' WHERE id = 'c21c5710-e2ff-4548-827b-0b26f83d3119';

-- Update the address
UPDATE addresses SET street_address_1 = 'Bldg 2150 Liggett Ave', street_address_2 = 'nan', city = 'JB Lewis-McChord', state = 'WA', postal_code = '98433' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c21c5710-e2ff-4548-827b-0b26f83d3119');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JB McGuire-Dix-Lakehurst - USAF', gbloc = 'AGFM' WHERE id = '10052404-6cf0-47da-b5bc-42b03e02f904';

-- Update the address
UPDATE addresses SET street_address_1 = '1706 Vandenberg Ave', street_address_2 = 'Bldg 1702, McGuire AFB', city = 'Trenton', state = 'NJ', postal_code = '08641' WHERE id = (SELECT address_id FROM transportation_offices where id = '10052404-6cf0-47da-b5bc-42b03e02f904');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JBSA Fort Sam Houston - USAF', gbloc = 'HAFC' WHERE id = 'e49796ab-bbae-4759-8b18-e0bf64eae315';

-- Update the address
UPDATE addresses SET street_address_1 = '2484 Stanley Rd', street_address_2 = 'Bldg 4023 Room 207', city = 'JBSA Fort Sam Houston', state = 'TX', postal_code = '78234' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e49796ab-bbae-4759-8b18-e0bf64eae315');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JBSA Lackland AFB - USAF', gbloc = 'HAFC' WHERE id = '3456949e-8615-4949-8899-4faa34f3880b';

-- Update the address
UPDATE addresses SET street_address_1 = '1651 Stewart St', street_address_2 = 'Bldg 5616, Room 112', city = 'JBSA Lackland', state = 'TX', postal_code = '78236' WHERE id = (SELECT address_id FROM transportation_offices where id = '3456949e-8615-4949-8899-4faa34f3880b');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO JBSA Randolph AFB - USAF', gbloc = 'HAFC' WHERE id = '55c4e988-9534-4bfd-863f-831c5c3d9421';

-- Update the address
UPDATE addresses SET street_address_1 = '550 D St East', street_address_2 = 'Bldg 399, Room 205', city = 'JBSA Randolph AFB', state = 'TX', postal_code = '78150' WHERE id = (SELECT address_id FROM transportation_offices where id = '55c4e988-9534-4bfd-863f-831c5c3d9421');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Keesler AFB - USAF', gbloc = 'HAFC' WHERE id = '7fcd0f3d-49de-4445-bd1d-074ed1a29215';

-- Update the address
UPDATE addresses SET street_address_1 = '500 Fisher St', street_address_2 = '81 LRS/LGRDF, Room 114', city = 'Keesler AFB', state = 'MS', postal_code = '39534' WHERE id = (SELECT address_id FROM transportation_offices where id = '7fcd0f3d-49de-4445-bd1d-074ed1a29215');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Key West - USN', gbloc = 'CNNQ' WHERE id = 'd29f36e7-003c-44e9-84f5-d8045f63fb87';

-- Update the address
UPDATE addresses SET street_address_1 = '811 Sigsbee Rd', street_address_2 = 'Bldg V-4059, Room 1', city = 'Key West', state = 'FL', postal_code = '33040' WHERE id = (SELECT address_id FROM transportation_offices where id = 'd29f36e7-003c-44e9-84f5-d8045f63fb87');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Kirtland AFB - USAF', gbloc = 'KKFA' WHERE id = '136eab07-558d-4ef8-aed5-c094b21ff31a';

-- Update the address
UPDATE addresses SET street_address_1 = '20245 4th St', street_address_2 = 'nan', city = 'Albuquerque', state = 'NM', postal_code = '87117' WHERE id = (SELECT address_id FROM transportation_offices where id = '136eab07-558d-4ef8-aed5-c094b21ff31a');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Laughlin AFB - USAF', gbloc = 'HAFC' WHERE id = 'c75f4700-c526-44d2-9f7e-f7d49beb2d56';

-- Update the address
UPDATE addresses SET street_address_1 = '427 Liberty Dr', street_address_2 = 'Bldg 77/246, Room 110', city = 'Laughlin AFB', state = 'TX', postal_code = '78843' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c75f4700-c526-44d2-9f7e-f7d49beb2d56');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Little Rock AFB - USAF', gbloc = 'HAFC' WHERE id = '7676dbfc-6719-4b6f-884a-036d1ce22454';

-- Update the address
UPDATE addresses SET street_address_1 = '1500 Vandenberg Blvd', street_address_2 = 'Bldg 1255, Suite 104', city = 'Little Rock AFB', state = 'AR', postal_code = '72099' WHERE id = (SELECT address_id FROM transportation_offices where id = '7676dbfc-6719-4b6f-884a-036d1ce22454');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Los Angeles SFB - USAF', gbloc = 'KKFA' WHERE id = 'ca6234a4-ed56-4094-a39c-738802798c6b';

-- Update the address
UPDATE addresses SET street_address_1 = '483 N. Aviation Blvd', street_address_2 = 'nan', city = 'Los Angeles SFB', state = 'CA', postal_code = '90245' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ca6234a4-ed56-4094-a39c-738802798c6b');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO LRC Dugway Proving Ground - USA', gbloc = 'KKFA' WHERE id = 'c8d5f80f-07c9-4d5f-9479-90b02e48c1d1';

-- Update the address
UPDATE addresses SET street_address_1 = 'Doolitte Ave', street_address_2 = 'Bldg 5464', city = 'Dugway', state = 'UT', postal_code = '84022' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c8d5f80f-07c9-4d5f-9479-90b02e48c1d1');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO LRC Redstone Arsenal - USA', gbloc = 'CNNQ' WHERE id = '5cdb638c-2649-45e3-b8d7-1b5ff5040228';

-- Update the address
UPDATE addresses SET street_address_1 = '3433 Snooper Rd', street_address_2 = 'Bldg 3433', city = 'Redstone Arsenal', state = 'AL', postal_code = '35898' WHERE id = (SELECT address_id FROM transportation_offices where id = '5cdb638c-2649-45e3-b8d7-1b5ff5040228');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Luke AFB - USAF', gbloc = 'KKFA' WHERE id = '3210a533-19b8-4805-a564-7eb452afce10';

-- Update the address
UPDATE addresses SET street_address_1 = '7383 N Litchfield Rd', street_address_2 = 'Room 1122', city = 'Luke AFB', state = 'AZ', postal_code = '85309' WHERE id = (SELECT address_id FROM transportation_offices where id = '3210a533-19b8-4805-a564-7eb452afce10');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO MacDill AFB - USAF', gbloc = 'HAFC' WHERE id = '21b421c5-c1cb-43ff-8b3b-65dabe6f39cc';

-- Update the address
UPDATE addresses SET street_address_1 = '8307 Cypress Stand Dr', street_address_2 = 'Bldg 49', city = 'MacDill AFB', state = 'FL', postal_code = '33621' WHERE id = (SELECT address_id FROM transportation_offices where id = '21b421c5-c1cb-43ff-8b3b-65dabe6f39cc');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Malmstrom AFB - USAF', gbloc = 'KKFA' WHERE id = '9804dbaa-71b0-4170-8ae2-fcba1356fb7f';

-- Update the address
UPDATE addresses SET street_address_1 = '7911 Goddard Dr', street_address_2 = '341st LRS/LGRT', city = 'Malmstrom AFB', state = 'MT', postal_code = '59402' WHERE id = (SELECT address_id FROM transportation_offices where id = '9804dbaa-71b0-4170-8ae2-fcba1356fb7f');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Maxwell AFB - USAF', gbloc = 'HAFC' WHERE id = '9215499f-4ac2-488f-b91d-6f14ec9d9160';

-- Update the address
UPDATE addresses SET street_address_1 = '50 LeMay Plaza S.', street_address_2 = 'Bldg 804', city = 'Maxwell AFB', state = 'AL', postal_code = '36112' WHERE id = (SELECT address_id FROM transportation_offices where id = '9215499f-4ac2-488f-b91d-6f14ec9d9160');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO McAlester - USA', gbloc = 'HAFC' WHERE id = 'd1359c20-c762-4b04-9ed6-fd2b9060615b';

-- Update the address
UPDATE addresses SET street_address_1 = '1 C Tree Rd', street_address_2 = 'Bldg 31', city = 'McAlester', state = 'OK', postal_code = '74501' WHERE id = (SELECT address_id FROM transportation_offices where id = 'd1359c20-c762-4b04-9ed6-fd2b9060615b');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO McChord Field - USA', gbloc = 'JEAT' WHERE id = '95abaeaa-452f-4fe0-9264-960cd2a15ccd';

-- Update the address
UPDATE addresses SET street_address_1 = '100 Colonel Joe Jackson Blvd', street_address_2 = 'Bldg 100 Custer Service Mall 1st Floor Room M6', city = 'JB Lewis-McChord', state = 'WA', postal_code = '98438' WHERE id = (SELECT address_id FROM transportation_offices where id = '95abaeaa-452f-4fe0-9264-960cd2a15ccd');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO McConnell AFB - USAF', gbloc = 'KKFA' WHERE id = '09932c15-8e61-47b6-83aa-402499d4366c';

-- Update the address
UPDATE addresses SET street_address_1 = '53476 Wichita St', street_address_2 = 'Room 340', city = 'McConnell AFB', state = 'KS', postal_code = '67221' WHERE id = (SELECT address_id FROM transportation_offices where id = '09932c15-8e61-47b6-83aa-402499d4366c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Minot AFB - USAF', gbloc = 'KKFA' WHERE id = '1ff3fe41-1b44-4bb4-8d5b-99d4ba5c2218';

-- Update the address
UPDATE addresses SET street_address_1 = '341 Bomber Blvd', street_address_2 = 'Bldg 527', city = 'Minot AFB', state = 'ND', postal_code = '58703' WHERE id = (SELECT address_id FROM transportation_offices where id = '1ff3fe41-1b44-4bb4-8d5b-99d4ba5c2218');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Moody AFB - USAF', gbloc = 'HAFC' WHERE id = 'a60b77a1-9f5f-48a8-80d6-579a7072bddb';

-- Update the address
UPDATE addresses SET street_address_1 = '23 Flying Tiger Way', street_address_2 = 'Bldg 182, Suite 179', city = 'Moody AFB', state = 'GA', postal_code = '31699' WHERE id = (SELECT address_id FROM transportation_offices where id = 'a60b77a1-9f5f-48a8-80d6-579a7072bddb');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Mountain Home AFB - USAF', gbloc = 'KKFA' WHERE id = '7ef50df6-f30e-4d82-9ca6-3130940315db';

-- Update the address
UPDATE addresses SET street_address_1 = '366 Gunfighter Ave', street_address_2 = 'Bldg 1132, Suite 242', city = 'Mountain Home AFB', state = 'ID', postal_code = '83648' WHERE id = (SELECT address_id FROM transportation_offices where id = '7ef50df6-f30e-4d82-9ca6-3130940315db');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAS Fallon - USN', gbloc = 'LKNQ' WHERE id = 'e665a46c-059e-4834-b7df-6e973747a92e';

-- Update the address
UPDATE addresses SET street_address_1 = '4755 Pasture Rd', street_address_2 = 'Bldg 66', city = 'Fallon', state = 'NV', postal_code = '89496' WHERE id = (SELECT address_id FROM transportation_offices where id = 'e665a46c-059e-4834-b7df-6e973747a92e');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAS JRB Fort Worth - USN', gbloc = 'HAFC' WHERE id = 'f69f315c-942a-4ef3-9427-7fee7883ce73';

-- Update the address
UPDATE addresses SET street_address_1 = '1330 Military Parkway', street_address_2 = 'nan', city = 'NAS JRB Fort Worth', state = 'TX', postal_code = '76127' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f69f315c-942a-4ef3-9427-7fee7883ce73');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAS Lemoore - USN', gbloc = 'LKNQ' WHERE id = '1039d189-39ba-47d4-8ed7-c96304576862';

-- Update the address
UPDATE addresses SET street_address_1 = '730 Enterprise Ave', street_address_2 = 'Wing 2, Room 5', city = 'NAS Lemoore', state = 'CA', postal_code = '93246' WHERE id = (SELECT address_id FROM transportation_offices where id = '1039d189-39ba-47d4-8ed7-c96304576862');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAS Pensacola - USN', gbloc = 'CNNQ' WHERE id = '2581f0aa-bc31-4c89-92cd-9a1843b49e59';

-- Update the address
UPDATE addresses SET street_address_1 = '121 Cuddahy St', street_address_2 = 'Bldg 680, Suite C', city = 'NAS Pensacola', state = 'FL', postal_code = '32508' WHERE id = (SELECT address_id FROM transportation_offices where id = '2581f0aa-bc31-4c89-92cd-9a1843b49e59');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Naval Postgraduate School Monterey - USN', gbloc = 'LKNQ' WHERE id = '2a0fb7ab-5a57-4450-b782-e7a58713bccb';

-- Update the address
UPDATE addresses SET street_address_1 = '1281 B Leahy Rd', street_address_2 = 'nan', city = 'Monterey', state = 'CA', postal_code = '93940' WHERE id = (SELECT address_id FROM transportation_offices where id = '2a0fb7ab-5a57-4450-b782-e7a58713bccb');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Naval Station Newport - USN', gbloc = 'AGFM' WHERE id = 'afafec09-7a91-4a7e-981d-3601f700ebbf';

-- Update the address
UPDATE addresses SET street_address_1 = '690 Peary St', street_address_2 = 'nan', city = 'Naval Station Newport', state = 'RI', postal_code = '02841' WHERE id = (SELECT address_id FROM transportation_offices where id = 'afafec09-7a91-4a7e-981d-3601f700ebbf');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Naval Training Center Great Lakes - USN', gbloc = 'BGNC' WHERE id = 'c65f1903-9f40-42e8-9720-f4e804702817';

-- Update the address
UPDATE addresses SET street_address_1 = '1710 B Cavin Dr', street_address_2 = 'Bldg 8100', city = 'Great Lakes', state = 'IL', postal_code = '60088' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c65f1903-9f40-42e8-9720-f4e804702817');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAVSUP Bethesda - USN', gbloc = 'BGAC' WHERE id = '0a58af30-a939-46a2-9b09-1fc40c5f0011';

-- Update the address
UPDATE addresses SET street_address_1 = '8901 Wisconsin Ave', street_address_2 = 'Bldg 17, Room 3D', city = 'Bethesda', state = 'MD', postal_code = '20889' WHERE id = (SELECT address_id FROM transportation_offices where id = '0a58af30-a939-46a2-9b09-1fc40c5f0011');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAVSUP FLC PS Everett - USN', gbloc = 'JEAT' WHERE id = '880fff01-d4be-4317-92f1-a3fab7ab1149';

-- Update the address
UPDATE addresses SET street_address_1 = '200 West Marine View Dr', street_address_2 = 'Bldg 2200', city = 'Everett', state = 'WA', postal_code = '98207' WHERE id = (SELECT address_id FROM transportation_offices where id = '880fff01-d4be-4317-92f1-a3fab7ab1149');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAVSUP FLC PS Whidbey Island - USN', gbloc = 'JEAT' WHERE id = 'f133c300-16af-4381-a1d7-a34edb094103';

-- Update the address
UPDATE addresses SET street_address_1 = '3675 W Lexington Dr', street_address_2 = 'BLDG 2556 Room 151', city = 'Oak Harbor', state = 'WA', postal_code = '98278' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f133c300-16af-4381-a1d7-a34edb094103');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAVSUP FLC Puget Sound - USN', gbloc = 'JEAT' WHERE id = 'affb700e-7e76-4fcc-a143-2c4ea4b0c480';

-- Update the address
UPDATE addresses SET street_address_1 = '2720 Ohio St', street_address_2 = 'Bangor Plaza Housing Office', city = 'Silverdale', state = 'WA', postal_code = '98315' WHERE id = (SELECT address_id FROM transportation_offices where id = 'affb700e-7e76-4fcc-a143-2c4ea4b0c480');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAVSUP FLC Seal Beach - USN', gbloc = 'LKNQ' WHERE id = 'a0de70b3-e7e9-4e47-8c7e-3f9f7ba4c5ab';

-- Update the address
UPDATE addresses SET street_address_1 = '800 Seal Beach Blvd', street_address_2 = 'Bldg 239', city = 'Seal Beach', state = 'CA', postal_code = '90740' WHERE id = (SELECT address_id FROM transportation_offices where id = 'a0de70b3-e7e9-4e47-8c7e-3f9f7ba4c5ab');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NAVSUP FLCN Annapolis - USN', gbloc = 'BGNC' WHERE id = '6eca781d-7b97-4893-afbe-2048c1629007';

-- Update the address
UPDATE addresses SET street_address_1 = '181 Wainwright Rd', street_address_2 = 'nan', city = 'Annapolis', state = 'MD', postal_code = '21402' WHERE id = (SELECT address_id FROM transportation_offices where id = '6eca781d-7b97-4893-afbe-2048c1629007');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Navy Submarine Base New London - USN', gbloc = 'AGFM' WHERE id = '5eb485ae-fb9c-4c90-80e4-6231158797df';

-- Update the address
UPDATE addresses SET street_address_1 = 'Grayling Ave', street_address_2 = 'Bldg 84, Room 206', city = 'Groton', state = 'CT', postal_code = '06349' WHERE id = (SELECT address_id FROM transportation_offices where id = '5eb485ae-fb9c-4c90-80e4-6231158797df');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NCBC Gulfport - USN', gbloc = 'HAFC' WHERE id = '3ea05911-86f9-4cf1-8f8e-b5e46ed36d4c';

-- Update the address
UPDATE addresses SET street_address_1 = '511 N Brown Ave', street_address_2 = 'Bldg 437', city = 'Gulfport', state = 'MS', postal_code = '39501' WHERE id = (SELECT address_id FROM transportation_offices where id = '3ea05911-86f9-4cf1-8f8e-b5e46ed36d4c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Nellis AFB - USAF', gbloc = 'KKFA' WHERE id = '91895530-82d3-4e06-9343-07710f3cdd29';

-- Update the address
UPDATE addresses SET street_address_1 = '4420 Grissom Ave', street_address_2 = 'Suite 203', city = 'Nellis AFB', state = 'NV', postal_code = '89191' WHERE id = (SELECT address_id FROM transportation_offices where id = '91895530-82d3-4e06-9343-07710f3cdd29');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NSA Mid-South Millington - USN', gbloc = 'CNNQ' WHERE id = 'c69bb3a5-7c7e-40bc-978f-e341f091ac52';

-- Update the address
UPDATE addresses SET street_address_1 = '766 Intrepid St', street_address_2 = 'Community Services Bldg 456, 1st Floor Left Wing', city = 'Millington', state = 'TN', postal_code = '38054' WHERE id = (SELECT address_id FROM transportation_offices where id = 'c69bb3a5-7c7e-40bc-978f-e341f091ac52');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO NSA Saratoga Springs - USN', gbloc = 'AGFM' WHERE id = '559eb724-7577-4e4c-830f-e55cbb030e06';

-- Update the address
UPDATE addresses SET street_address_1 = '19 J. F. King Dr', street_address_2 = 'Bldg 102', city = 'Saratoga Springs', state = 'NY', postal_code = '12866' WHERE id = (SELECT address_id FROM transportation_offices where id = '559eb724-7577-4e4c-830f-e55cbb030e06');

-- Insert the address
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('ec9c2d10-7934-432a-8c1b-a87cefeb1ed1', '110 Vernon Ave', 'Bldg 386', 'Panama City', 'FL', '32407', now(), now());

-- Insert the TO
INSERT INTO transportation_offices (id, address_id, gbloc, name, created_at, updated_at, latitude, longitude) VALUES ('57cf1e81-8113-4a52-bc50-3cb8902c2efd', 'ec9c2d10-7934-432a-8c1b-a87cefeb1ed1', 'HAFC', 'PPPO NSWC Panama City Division - USN', now(), now(), 0, 0);

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Offutt AFB - USAF', gbloc = 'KKFA' WHERE id = '78443187-1d59-4e10-82c9-ee7a0e062730';

-- Update the address
UPDATE addresses SET street_address_1 = '106 Peacekeeper Dr', street_address_2 = 'Suite 2N3', city = 'Offutt AFB', state = 'NE', postal_code = '68113' WHERE id = (SELECT address_id FROM transportation_offices where id = '78443187-1d59-4e10-82c9-ee7a0e062730');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Patrick SFB - USAF', gbloc = 'HAFC' WHERE id = '261da328-1652-4d33-b590-9e62f1137203';

-- Update the address
UPDATE addresses SET street_address_1 = '972 South Patrick Dr', street_address_2 = 'Bldg 821', city = 'Patrick SFB', state = 'FL', postal_code = '32925' WHERE id = (SELECT address_id FROM transportation_offices where id = '261da328-1652-4d33-b590-9e62f1137203');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Peterson SFB - USAF', gbloc = 'KKFA' WHERE id = 'cc107598-3d72-4679-a4aa-c28d1fd2a016';

-- Update the address
UPDATE addresses SET street_address_1 = '135 Dover St', street_address_2 = 'Bldg 350, Room 1046', city = 'Peterson SFB', state = 'CO', postal_code = '80914' WHERE id = (SELECT address_id FROM transportation_offices where id = 'cc107598-3d72-4679-a4aa-c28d1fd2a016');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Pine Bluff Arsenal - USA', gbloc = 'HAFC' WHERE id = '65363850-3cb1-44c7-a7c1-191b28a53479';

-- Update the address
UPDATE addresses SET street_address_1 = '10020 Kabrich Cir', street_address_2 = 'Bldg 11-080', city = 'Pine Bluff', state = 'AR', postal_code = '71602' WHERE id = (SELECT address_id FROM transportation_offices where id = '65363850-3cb1-44c7-a7c1-191b28a53479');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Portsmouth Navy Shipyard - USN', gbloc = 'AGFM' WHERE id = '30ad1395-fc3b-4a16-839d-865b19898f8d';

-- Update the address
UPDATE addresses SET street_address_1 = 'PPPO Portsmouth Navy Ship Yard', street_address_2 = 'Bldg H10, Room 6', city = 'Kittery', state = 'ME', postal_code = '03904' WHERE id = (SELECT address_id FROM transportation_offices where id = '30ad1395-fc3b-4a16-839d-865b19898f8d');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Presidio of Monterey - USA', gbloc = 'LKNQ' WHERE id = 'b6e7afe2-a58c-45cc-b65b-3abda89a1ed6';

-- Update the address
UPDATE addresses SET street_address_1 = '1712 Private Bolio Rd', street_address_2 = 'Bldg 517', city = 'Monterey', state = 'CA', postal_code = '93944' WHERE id = (SELECT address_id FROM transportation_offices where id = 'b6e7afe2-a58c-45cc-b65b-3abda89a1ed6');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Red River Army Depot - USA', gbloc = 'HAFC' WHERE id = '6192fd68-d5f2-4f06-9cb0-f1a1c8f00b93';

-- Update the address
UPDATE addresses SET street_address_1 = '100 James Carlow Dr', street_address_2 = 'Bldg 431N', city = 'Texarkana', state = 'TX', postal_code = '75507' WHERE id = (SELECT address_id FROM transportation_offices where id = '6192fd68-d5f2-4f06-9cb0-f1a1c8f00b93');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Robins AFB - USAF', gbloc = 'HAFC' WHERE id = '540e2057-7a6a-4a39-8ff5-2e472c4fcd25';

-- Update the address
UPDATE addresses SET street_address_1 = '375 Perry St', street_address_2 = 'Bldg 914, Suite 208', city = 'Robins AFB', state = 'GA', postal_code = '31098' WHERE id = (SELECT address_id FROM transportation_offices where id = '540e2057-7a6a-4a39-8ff5-2e472c4fcd25');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Rock Island Arsenal - USA', gbloc = 'AGFM' WHERE id = '9e83a154-ae38-47a2-98da-52b38f4a87a1';

-- Update the address
UPDATE addresses SET street_address_1 = '1 Rock Island Arsenal', street_address_2 = 'ASPA-LRI Bldg 102 1st Fl NE Wing', city = 'Rock Island', state = 'IL', postal_code = '61299' WHERE id = (SELECT address_id FROM transportation_offices where id = '9e83a154-ae38-47a2-98da-52b38f4a87a1');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Rome - USAF', gbloc = 'AGFM' WHERE id = '767e8893-0063-404c-b9c5-df8f4a12cb70';

-- Update the address
UPDATE addresses SET street_address_1 = '26 Electronic Parkway', street_address_2 = 'Bldg 2, Room C-250', city = 'Rome', state = 'NY', postal_code = '13441' WHERE id = (SELECT address_id FROM transportation_offices where id = '767e8893-0063-404c-b9c5-df8f4a12cb70');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO San Diego - USMC', gbloc = 'LKNQ' WHERE id = '7e6b9019-5493-40a4-9dcd-c83fb4f77961';

-- Update the address
UPDATE addresses SET street_address_1 = '4100 Hochmuth Ave', street_address_2 = 'nan', city = 'San Diego', state = 'CA', postal_code = '92140' WHERE id = (SELECT address_id FROM transportation_offices where id = '7e6b9019-5493-40a4-9dcd-c83fb4f77961');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Scott AFB - USAF', gbloc = 'AGFM' WHERE id = '0931a9dc-c1fd-444a-b138-6e1986b1714c';

-- Update the address
UPDATE addresses SET street_address_1 = '215 Heritage Dr', street_address_2 = 'Bldg P-8, Room D-100', city = 'Scott AFB', state = 'IL', postal_code = '62225' WHERE id = (SELECT address_id FROM transportation_offices where id = '0931a9dc-c1fd-444a-b138-6e1986b1714c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Sector New York  - USCG', gbloc = 'BGAC' WHERE id = '19bafaf1-8e6f-492d-b6ac-6eacc1e5b64c';

-- Update the address
UPDATE addresses SET street_address_1 = '311 Battery Rd', street_address_2 = 'nan', city = 'Staten Island', state = 'NY', postal_code = '10305' WHERE id = (SELECT address_id FROM transportation_offices where id = '19bafaf1-8e6f-492d-b6ac-6eacc1e5b64c');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Seymour Johnson AFB - USAF', gbloc = 'AGFM' WHERE id = 'f760ee76-d386-47d8-9df1-aff8c347fa96';

-- Update the address
UPDATE addresses SET street_address_1 = '1280 Humphreys St', street_address_2 = 'nan', city = 'Seymour Johnson AFB', state = 'NC', postal_code = '27531' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f760ee76-d386-47d8-9df1-aff8c347fa96');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Shaw AFB - USAF', gbloc = 'AGFM' WHERE id = 'bf32bb9f-f0fd-4c3f-905f-0d88b3798a81';

-- Update the address
UPDATE addresses SET street_address_1 = '524 Shaw Dr', street_address_2 = 'nan', city = 'Shaw AFB', state = 'SC', postal_code = '29152' WHERE id = (SELECT address_id FROM transportation_offices where id = 'bf32bb9f-f0fd-4c3f-905f-0d88b3798a81');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Sheppard AFB - USAF', gbloc = 'HAFC' WHERE id = '6dc9c87d-37af-4157-a1ce-af45bb954eba';

-- Update the address
UPDATE addresses SET street_address_1 = '426 5th Ave', street_address_2 = 'Bldg 402, Suite 10', city = 'Sheppard AFB', state = 'TX', postal_code = '76311' WHERE id = (SELECT address_id FROM transportation_offices where id = '6dc9c87d-37af-4157-a1ce-af45bb954eba');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Surface Forces Logistics Command Baltimore - USCG', gbloc = 'BGAC' WHERE id = '0a048293-15c4-4036-8915-dd4b9d3ef2de';

-- Update the address
UPDATE addresses SET street_address_1 = 'ATTN: Transportion Office ', street_address_2 = '2401 Hawkins Point Rd', city = 'Baltimore', state = 'MD', postal_code = '21226' WHERE id = (SELECT address_id FROM transportation_offices where id = '0a048293-15c4-4036-8915-dd4b9d3ef2de');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Tinker AFB - USAF', gbloc = 'HAFC' WHERE id = '7876373d-57e4-4cde-b11f-c26a8feee9e8';

-- Update the address
UPDATE addresses SET street_address_1 = '7330 Century Blvd', street_address_2 = 'Bldg 469', city = 'Tinker AFB', state = 'OK', postal_code = '73145' WHERE id = (SELECT address_id FROM transportation_offices where id = '7876373d-57e4-4cde-b11f-c26a8feee9e8');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Tobyhanna Army Depot - USA', gbloc = 'AGFM' WHERE id = '46898e12-8657-4ece-bb89-9a9e94815db9';

-- Update the address
UPDATE addresses SET street_address_1 = '11 Hap Arnold Blvd', street_address_2 = 'nan', city = 'Coolbaugh Township', state = 'PA', postal_code = '18466' WHERE id = (SELECT address_id FROM transportation_offices where id = '46898e12-8657-4ece-bb89-9a9e94815db9');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Training Center Cape May - USCG', gbloc = 'AGFM' WHERE id = '7ac0f374-b97a-4b98-9878-eabef89adff9';

-- Update the address
UPDATE addresses SET street_address_1 = 'ATTN: Transportation Office  ', street_address_2 = '1 Munro Ave', city = 'Cape May', state = 'NJ', postal_code = '08204' WHERE id = (SELECT address_id FROM transportation_offices where id = '7ac0f374-b97a-4b98-9878-eabef89adff9');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Training Center Petaluma - USCG', gbloc = 'LHNQ' WHERE id = 'f54d8b95-6ee8-4ffa-bf79-67400ae09aa2';

-- Update the address
UPDATE addresses SET street_address_1 = '599 Tomales Rd', street_address_2 = 'Bldg 150', city = 'Petaluma', state = 'CA', postal_code = '94952' WHERE id = (SELECT address_id FROM transportation_offices where id = 'f54d8b95-6ee8-4ffa-bf79-67400ae09aa2');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Training Center Yorktown - USCG', gbloc = 'AGFM' WHERE id = '4ed2762d-73bc-4c62-bea9-725c5c64cb62';

-- Update the address
UPDATE addresses SET street_address_1 = 'ATTN: Transportation Office ', street_address_2 = '1 US Coast Guard Training Center', city = 'Yorktown', state = 'VA', postal_code = '23690' WHERE id = (SELECT address_id FROM transportation_offices where id = '4ed2762d-73bc-4c62-bea9-725c5c64cb62');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Travis AFB - USAF', gbloc = 'KKFA' WHERE id = 'a038e200-8db4-499f-b1a3-2c15f6e97614';

-- Update the address
UPDATE addresses SET street_address_1 = '540 Airlift Dr', street_address_2 = 'nan', city = 'Travis AFB', state = 'CA', postal_code = '94535' WHERE id = (SELECT address_id FROM transportation_offices where id = 'a038e200-8db4-499f-b1a3-2c15f6e97614');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Tyndall AFB - USAF', gbloc = 'HAFC' WHERE id = '4a2e6595-46e0-477c-b037-dfe042283ebe';

-- Update the address
UPDATE addresses SET street_address_1 = '445 Suwanne Rd', street_address_2 = 'Bldg 662, Suite 152', city = 'Tyndall AFB', state = 'FL', postal_code = '32403' WHERE id = (SELECT address_id FROM transportation_offices where id = '4a2e6595-46e0-477c-b037-dfe042283ebe');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO USAF Academy - USAF', gbloc = 'KKFA' WHERE id = '0452f996-f3e0-4e32-a48b-e5249c6e3d78';

-- Update the address
UPDATE addresses SET street_address_1 = '5136 Eagle Dr', street_address_2 = 'nan', city = 'U.S. Air Force Academy', state = 'CO', postal_code = '80840' WHERE id = (SELECT address_id FROM transportation_offices where id = '0452f996-f3e0-4e32-a48b-e5249c6e3d78');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Vance AFB - USAF', gbloc = 'HAFC' WHERE id = '2be0d1b8-68f0-4aee-b86b-5babbd5d49af';

-- Update the address
UPDATE addresses SET street_address_1 = '400 Young Rd', street_address_2 = 'Bldg 200, Suite 202', city = 'Vance AFB', state = 'OK', postal_code = '73705' WHERE id = (SELECT address_id FROM transportation_offices where id = '2be0d1b8-68f0-4aee-b86b-5babbd5d49af');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Vandenberg SFB - USAF', gbloc = 'KKFA' WHERE id = '32b29875-474a-4983-98c2-c02694d10724';

-- Update the address
UPDATE addresses SET street_address_1 = '1221 California Blvd', street_address_2 = 'Bldg 11777', city = 'Vandenberg SFB', state = 'CA', postal_code = '93437' WHERE id = (SELECT address_id FROM transportation_offices where id = '32b29875-474a-4983-98c2-c02694d10724');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Ventura County Site - USN', gbloc = 'LKNQ' WHERE id = '8b78ce40-0b34-413d-98af-c51d440e7a4d';

-- Update the address
UPDATE addresses SET street_address_1 = '1000 23rd Ave', street_address_2 = 'Bldg 1169, Room 1200', city = 'Port Hueneme', state = 'CA', postal_code = '93043' WHERE id = (SELECT address_id FROM transportation_offices where id = '8b78ce40-0b34-413d-98af-c51d440e7a4d');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO West Point Military Academy - USA', gbloc = 'BGAC' WHERE id = '46898e12-8657-4ece-bb89-9a9e94815db9';

-- Update the address
UPDATE addresses SET street_address_1 = 'ITO 626 Swift Rd', street_address_2 = 'nan', city = 'West Point', state = 'NY', postal_code = '10996' WHERE id = (SELECT address_id FROM transportation_offices where id = '46898e12-8657-4ece-bb89-9a9e94815db9');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO White Sands Missile Range - USA', gbloc = 'KKFA' WHERE id = 'ff6f1f7c-5309-4436-b9ff-e5dcac95b750';

-- Update the address
UPDATE addresses SET street_address_1 = '143 Crozier St', street_address_2 = 'Room 216', city = 'White Sands Missile Range', state = 'NM', postal_code = '88002' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ff6f1f7c-5309-4436-b9ff-e5dcac95b750');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Whiteman AFB - USAF', gbloc = 'KKFA' WHERE id = '658e2bce-b24c-4972-89c8-1676242bacdc';

-- Update the address
UPDATE addresses SET street_address_1 = '727 2nd St', street_address_2 = 'Suite 129', city = 'Whiteman AFB', state = 'MO', postal_code = '65305' WHERE id = (SELECT address_id FROM transportation_offices where id = '658e2bce-b24c-4972-89c8-1676242bacdc');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Wright-Patterson AFB - USAF', gbloc = 'AGFM' WHERE id = '9ac9a242-d193-49b9-b24a-f2825452f737';

-- Update the address
UPDATE addresses SET street_address_1 = '1940 Allbrook Dr', street_address_2 = 'nan', city = 'Wright-Patterson AFB', state = 'OH', postal_code = '45433' WHERE id = (SELECT address_id FROM transportation_offices where id = '9ac9a242-d193-49b9-b24a-f2825452f737');

-- Update the TO
UPDATE transportation_offices SET name = 'PPPO Yuma Proving Ground - USA', gbloc = 'LKNQ' WHERE id = '7145e8fe-465b-44e5-a486-b893357148ef';

-- Update the address
UPDATE addresses SET street_address_1 = '301 C St', street_address_2 = 'Bldg 2710', city = 'Yuma', state = 'AZ', postal_code = '85365' WHERE id = (SELECT address_id FROM transportation_offices where id = '7145e8fe-465b-44e5-a486-b893357148ef');

-- Update the TO
UPDATE transportation_offices SET name = 'PPSO DMO Camp Lejeune - USMC ', gbloc = 'USMC' WHERE id = 'ccf50409-9d03-4cac-a931-580649f1647a';

-- Update the address
UPDATE addresses SET street_address_1 = 'Ash St', street_address_2 = 'Bldg 1011', city = 'Camp Lejeune', state = 'NC', postal_code = '28547' WHERE id = (SELECT address_id FROM transportation_offices where id = 'ccf50409-9d03-4cac-a931-580649f1647a');

-- Update the TO
UPDATE transportation_offices SET name = 'PPSO Miami - USA', gbloc = 'CLPK' WHERE id = '7f7cc97c-2f3c-4866-90fe-b335f5c8e042';

-- Update the address
UPDATE addresses SET street_address_1 = '9301 NW 33rd St', street_address_2 = 'Suite C1029', city = 'Miami', state = 'FL', postal_code = '33172' WHERE id = (SELECT address_id FROM transportation_offices where id = '7f7cc97c-2f3c-4866-90fe-b335f5c8e042');

-- Update the TO
UPDATE transportation_offices SET name = 'PPSO NAVSUP FLC China Lake - USN', gbloc = 'LKNQ' WHERE id = '7e50b5f0-1717-4067-95d5-a2adb41939c5';

-- Update the address
UPDATE addresses SET street_address_1 = '1 Administration Cir', street_address_2 = 'Bldg 1023', city = 'China Lake', state = 'CA', postal_code = '93555' WHERE id = (SELECT address_id FROM transportation_offices where id = '7e50b5f0-1717-4067-95d5-a2adb41939c5');

-- Update the TO
UPDATE transportation_offices SET name = 'USN PPM Processing (NAVY) - USN', gbloc = 'NAVY' WHERE id = '83d25227-e27d-4c56-8e71-1ff464f47128';

-- Update the address
UPDATE addresses SET street_address_1 = '1968 Gilbert St', street_address_2 = 'Suite 600', city = 'Norfolk', state = 'VA', postal_code = '23511' WHERE id = (SELECT address_id FROM transportation_offices where id = '83d25227-e27d-4c56-8e71-1ff464f47128');
