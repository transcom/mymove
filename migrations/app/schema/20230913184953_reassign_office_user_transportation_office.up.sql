-- Jira ticket: MB-17020
-- Reassign office users to new transportation offices

-- 'Army Sustainment Command DTA Warren MI' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '28f4f8ef-5a79-420f-9837-e558a04ba060';

-- 'Base Ketchikan' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26938'
	WHERE transportation_office_id = '4afd7912-5cb5-4a90-a85d-ec72b436380e';

-- 'Camp Smith Personal Property Office' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '8ee1c261-198f-4efd-9716-6b2ee3cb5e80';

-- 'CG BASE KETCHIKAN' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26939'
	WHERE transportation_office_id = '2d65fc57-ab6b-4965-b5b2-cde15f166cb3';

-- 'CG BASE KETCHIKAN, AK' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26940'
	WHERE transportation_office_id = '0b2545a6-bc74-4c35-b7fb-eea2647cbbb7';

-- 'CG BASE KODIAK, AK' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26941'
	WHERE transportation_office_id = 'b1ceb0a7-9457-4595-b61f-fb89ef668f1f';

-- 'CG BSU, KODIAK, AK' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26942'
	WHERE transportation_office_id = 'a617a56f-1e8c-4de3-bfce-81e4780361c2';

-- 'Charleston Naval Weapon Station' -> 'PPPO JB Charleston (Naval Weapon Station) - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '67281ea0-222a-41ea-9ec2-dc274a2513ea'
	WHERE transportation_office_id = 'ac6af9b2-fcb0-4b59-98f4-6a8071bd3cba';

-- 'CPPSO - Carlilse Barracks, PA' -> 'PPPO Carlisle Barracks - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'e37860be-c642-4037-af9a-8a1be690d8d7'
	WHERE transportation_office_id = '8736624a-09d6-4867-b712-2287b3df766a';

-- 'CPPSO, CARLISLE BARRACKS, PA' -> 'PPPO Carlisle Barracks - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'e37860be-c642-4037-af9a-8a1be690d8d7'
	WHERE transportation_office_id = '1950314f-9527-4141-b116-cdc1f857f960';

-- 'CPPSO - HD, FORT HOOD, TX' -> 'CPPSO Fort Cavazos (HBAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b'
	WHERE transportation_office_id = 'd8acfa3d-725f-46dc-a272-658d4681b97f';

-- 'Davis Monthan AFB' -> 'PPPO Davis Monthan AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '54156892-dff1-4657-8998-39ff4e3a259e'
	WHERE transportation_office_id = '3977e742-ac3b-458d-bed9-db1120786447';

-- 'Eielson AFB' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26943'
	WHERE transportation_office_id = '41ef1e1c-c257-48d3-8727-ba560ac6ac3d';

-- 'Ellsworth AFB' -> 'PPPO Ellsworth AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '07776e49-5100-4094-a7b9-5d9de4fa18d6'
	WHERE transportation_office_id = 'aad5143b-f9ee-4d56-999b-9a61e78bfd69';

-- 'FLEET LOGISTIC CENTER PEARL HARBOR' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff';

-- 'FORT BRAGG, NC' -> 'PPPO Fort Liberty - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'e3c44c50-ece0-4692-8aad-8c9600000000'
	WHERE transportation_office_id = '8eca374c-f5b1-4f88-8821-1d82068a0cbf';

-- 'JBER Travel Center-Elmendorf' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26944'
	WHERE transportation_office_id = '4522d141-87f1-4f1e-a111-466303c6ae14';

-- 'JBER Travel Center-Richardson' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26945'
	WHERE transportation_office_id = 'bc34e876-7f18-4401-ab91-507b0861a947';

-- 'JOINT BASE ELMENDORF-RICHARDSON, AK' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26946'
	WHERE transportation_office_id = '0dcf17dd-e06a-435f-91cf-ccef70af35e0';

-- 'JOINT PERS PROP SHIPPING OFFICE - MA' -> 'JPPSO - Mid Atlantic (BGAC) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a50'
	WHERE transportation_office_id = 'b97a217c-daac-4ce8-8a30-8c914f6812f1';

-- 'JPPSO: JOINT PERS PROP SHIPPING OFFICE - MA' -> 'JPPSO - Mid Atlantic (BGAC) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a51'
	WHERE transportation_office_id = '1b6d6f26-b0d9-4d3d-ab90-33816ced0c83';

-- 'JPPSO-MA: Regional Customer Service Office' -> 'JPPSO - Mid Atlantic (BGAC) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a52'
	WHERE transportation_office_id = '39c5c8a3-b758-49de-8606-588a8a67b149';

-- 'JPPSO SOUTHWEST' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '0509ae13-9216-41ed-a7e1-c521732e03ef';

-- 'Kaneohe Marine Corps Air Station Personal Property Office' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '8cb285cd-576e-4325-a02b-a2050cc559e8';

-- 'McChord AFB' -> 'PPPO JB Lewis-McChord (Fort Lewis) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'c21c5710-e2ff-4548-827b-0b26f83d3119'
	WHERE transportation_office_id = '72823c31-a622-4c82-9e8f-0a11e4e4654e';

-- 'MEPS' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '8f849395-2dfb-48ee-9bc5-704189a3b366';

-- 'NAVSUP FLC NORFOLK-CPPSO' -> 'CPPSO Norfolk (BGNC) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5f741385-0a34-4d05-9068-e1e2dd8dfefc'
	WHERE transportation_office_id = 'cc6ce760-c324-4b83-bed7-c24153fd074c';

-- 'NAVSUP FLC PPSO CHINA LAKE' -> 'PPSO NAVSUP FLC China Lake - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '7e50b5f0-1717-4067-95d5-a2adb41939c5'
	WHERE transportation_office_id = '20eff1d4-e190-4578-8d45-03910360f310';

-- 'NAVSUP FLC PUGET SOUND' -> 'PPPO NAVSUP FLC Puget Sound - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'affb700e-7e76-4fcc-a143-2c4ea4b0c480'
	WHERE transportation_office_id = 'eadd62ac-e17f-4d36-97e5-8cc1b40a28ac';

-- 'NAVSUP FLC PUGET SOUND' -> 'PPPO NAVSUP FLC Puget Sound - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'affb700e-7e76-4fcc-a143-2c4ea4b0c480'
	WHERE transportation_office_id = '2726d99c-eeaa-4c54-b24f-692ba0a78e2b';

-- 'NSA MidSouth Personal Property Office' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '04df1fa8-a0e9-4a1e-9d03-7e37947ce81d';

-- 'PACIFIC MISSILE RANGE FACILITY (KAUAI) - NAVY' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '3e1a2171-0f6a-4c87-9cab-cfa7c0bcecb3';

-- 'Personal Property Dept, Joint Base Pearl Harbor Hickam (JBPHH)' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '071b9dfe-039e-4e5b-b493-010aec575f0e';

-- 'Personal Property Office' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = 'd1fdd946-6f29-4eef-ad4f-54c213e84d9e';

-- 'Personal Property Processing Office Fort Greely, AK' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26947'
	WHERE transportation_office_id = 'dd2c98a6-303d-4596-86e8-b067a7deb1a2';

-- 'Personal Property Processing Office Fort Leavenworth' -> 'PPPO Fort Leavenworth - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'b2f76d56-6996-41a3-aef7-483524a643d1'
	WHERE transportation_office_id = '462d000b-3cc8-4c38-94d6-4e8b540affec';

-- 'Personal Property Processing Office Fort Leonard Wood' -> 'PPPO Fort Leonard Wood - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7'
	WHERE transportation_office_id = '631d129b-5d1c-4bed-b37f-d54c423832ef';

-- 'Personal Property Processing Office Fort Wainwright, AK' -> 'JPPSO - North West (JEAT) - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26948'
	WHERE transportation_office_id = '446aaf44-a5c8-4000-a0b8-6e5e421f62b0';

-- 'Personal Property Processing Office Lewis-Main' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '56f61173-214a-4498-9f76-39f22890aea4';

-- 'PPO Seal Beach NWS, CA' -> 'PPPO NAVSUP FLC Seal Beach - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'a0de70b3-e7e9-4e47-8c7e-3f9f7ba4c5ab'
	WHERE transportation_office_id = 'bc3e0b87-7a63-44be-b551-1510e8e24655';

-- 'PPPO FORT BENNING, GA' -> 'PPPO Fort Moore - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '5a9ed15c-ed78-47e5-8afd-7583f3cc660d'
	WHERE transportation_office_id = '524e0511-fea5-4167-964a-280f41708a97';

-- 'PPSO, FORT LEONARD WOOD, MO' -> 'PPPO Fort Leonard Wood - USA'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7'
	WHERE transportation_office_id = '9bcdf7fb-9bdc-4f1f-9d07-8b7520c9f3fd';

-- 'Schofield Barracks  / Fort Shafter Transportation office' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = 'b0d787ad-94f8-4bb6-8230-85bad755f07c';

-- 'Sector Humboldt Bay' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '15ca034d-89bb-4124-ad20-4b75d7f0c101';

-- 'Tooele Army Depot' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '41b95c70-f4a5-4e55-94e2-4d3ba44c8e38';

-- 'Transportation Personal Property' -> 'PPPO Scott AFB - USAF'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '0931a9dc-c1fd-444a-b138-6e1986b1714c'
	WHERE transportation_office_id = '4858d22c-ec59-4ec6-985f-73366bfa08c1';

-- 'USCG BASE ALAMEDA, CA' -> 'PPPO Base Alameda - USCG'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '3fc4b408-1197-430a-a96a-24a5a1685b45'
	WHERE transportation_office_id = '9c509a6f-e87c-4e1f-b04d-780ddaf4d340';

-- 'USCG Base Honolulu Transportation' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE office_users
	SET updated_at = now(), transportation_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
	WHERE transportation_office_id = '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58';
