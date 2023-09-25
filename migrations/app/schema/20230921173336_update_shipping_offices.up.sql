-- Jira: MB-17020
-- Updates shipping offices for a number of transportation offices

-- 'JPPSO SOUTHWEST' -> 'JPPSO - South West (LKNQ) - USN'
UPDATE transportation_offices
    SET shipping_office_id = '27002d34-e9ea-4ef5-a086-f23d07c4088c'
    WHERE shipping_office_id = '0509ae13-9216-41ed-a7e1-c521732e03ef';

-- 'CPPSO, CARLISLE BARRACKS, PA' -> 'PPPO Carlisle Barracks - USA'
UPDATE transportation_offices
    SET shipping_office_id = 'e37860be-c642-4037-af9a-8a1be690d8d7'
    WHERE shipping_office_id = '1950314f-9527-4141-b116-cdc1f857f960';

-- 'FORT BRAGG, NC' -> 'JPPSO - Mid Atlantic (BGAC) - USA'
UPDATE transportation_offices
    SET shipping_office_id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a50'
    WHERE shipping_office_id = '8eca374c-f5b1-4f88-8821-1d82068a0cbf';

-- 'PPSO, FORT LEONARD WOOD, MO' -> 'PPPO Fort Leonard Wood - USA'
UPDATE transportation_offices
    SET shipping_office_id = 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7'
    WHERE shipping_office_id = '9bcdf7fb-9bdc-4f1f-9d07-8b7520c9f3fd';

-- 'USCG BASE ALAMEDA, CA' -> 'PPPO Base Alameda - USCG'
UPDATE transportation_offices
    SET shipping_office_id = '3fc4b408-1197-430a-a96a-24a5a1685b45'
    WHERE shipping_office_id = '9c509a6f-e87c-4e1f-b04d-780ddaf4d340';

-- 'JOINT PERS PROP SHIPPING OFFICE - MA' -> 'JPPSO - Mid Atlantic (BGAC) - USA'
UPDATE transportation_offices
    SET shipping_office_id = '8e25ccc1-7891-4146-a9d0-cd0d48b59a50'
    WHERE shipping_office_id = 'b97a217c-daac-4ce8-8a30-8c914f6812f1';

-- 'NAVSUP FLC NORFOLK-CPPSO' -> 'CPPSO Norfolk (BGNC) - USN'
UPDATE transportation_offices
    SET shipping_office_id = '5f741385-0a34-4d05-9068-e1e2dd8dfefc'
    WHERE shipping_office_id = 'cc6ce760-c324-4b83-bed7-c24153fd074c';

-- 'CPPSO - HD, FORT HOOD, TX' -> 'CPPSO Fort Cavazos (HBAT) - USA'
UPDATE transportation_offices
    SET shipping_office_id = 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b'
    WHERE shipping_office_id = 'd8acfa3d-725f-46dc-a272-658d4681b97f';

-- 'NAVSUP FLC PUGET SOUND' -> 'PPPO NAVSUP FLC Puget Sound - USN'
UPDATE transportation_offices
    SET shipping_office_id = 'affb700e-7e76-4fcc-a143-2c4ea4b0c480'
    WHERE shipping_office_id = 'eadd62ac-e17f-4d36-97e5-8cc1b40a28ac';
