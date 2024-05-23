update postal_code_to_gblocs set gbloc = 'BGAC' where gbloc = 'BKAS';


INSERT INTO addresses
    (id, street_address_1,  city, state, postal_code, created_at, updated_at, country)
    VALUES ('f933c50f-6625-4991-8c81-705a222840c6', '3376 Albacore Alley', 'San Diego', 'CA', '92136', now(), now(), 'United States');

insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at, provides_ppm_closeout)
values
    ('e4a02d40-2ad9-44c5-a357-202d4ff0b51d', 'PPPO NAVSUP FLC San Diego - USN', 'LKNQ', 'f933c50f-6625-4991-8c81-705a222840c6','32.67540','-117.12142', now(), now(), TRUE);

insert into office_phone_lines
    (id, transportation_office_id, number, created_at, updated_at)
values
    ('3e692f01-182d-4a7d-958d-6418e7335dd9', 'e4a02d40-2ad9-44c5-a357-202d4ff0b51d', '855-444-6683', now(), now());

insert into office_emails
    (id, transportation_office_id, email, created_at, updated_at)
values
    ('3e692f01-182d-4a7d-958d-6418e7335dd9', 'e4a02d40-2ad9-44c5-a357-202d4ff0b51d', 'jppso_SW_counseling@us.navy.mil', now(), now());



INSERT INTO addresses
    (id, street_address_1,  city, state, postal_code, created_at, updated_at, country)
    VALUES ('6aa77b74-41a7-4a4c-ab29-986f3263495a', '626 Swift Road', 'West Point', 'NY', '10996', now(), now(), 'United States');

insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at, provides_ppm_closeout)
values
    ('dd043073-4f1b-460f-8f8c-74403619dbaa', 'PPPO West Point/ USMA - USA', 'BGAC', '6aa77b74-41a7-4a4c-ab29-986f3263495a','41.39400','-73.97232', now(), now(), TRUE);
insert into office_phone_lines
    (id, transportation_office_id, number, created_at, updated_at)
values
    ('1adffbea-01d7-462f-bbfd-ba155e8a0844', 'dd043073-4f1b-460f-8f8c-74403619dbaa', '845-938-5911', now(), now());

insert into office_emails
    (id, transportation_office_id, email, created_at, updated_at)
values
    ('8276c5e4-461f-4ef6-a72f-0d76f7e10194', 'dd043073-4f1b-460f-8f8c-74403619dbaa', 'usarmy.jblm.404-afsb-lrc.list.west-point-transportation-ppo@army.mil', now(), now());


INSERT INTO addresses
    (id, street_address_1,  city, state, postal_code, created_at, updated_at, country)
    VALUES ('09058d36-2966-496a-aaf5-55c024404396', '15610 SW 117TH AVE', 'Miami', 'FL', '33177-1630', now(), now(), 'United States');

insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at, provides_ppm_closeout)
values
    ('4f10d0f5-6017-4de2-8cfb-ee9252e492d5', 'PPPO USAG Miami - USA', 'CLPK', '09058d36-2966-496a-aaf5-55c024404396','25.59788','-80.40353', now(), now(), TRUE);

insert into office_phone_lines
    (id, transportation_office_id, number, created_at, updated_at)
values
    ('bfae8dd6-3ce4-4310-a315-eeda8420f4a4', '4f10d0f5-6017-4de2-8cfb-ee9252e492d5', '1-305-216-8037', now(), now());

insert into office_emails
    (id, transportation_office_id, email, created_at, updated_at)
values
    ('afb6648e-da9d-4e6a-9322-e0aa3b99caf3', '4f10d0f5-6017-4de2-8cfb-ee9252e492d5', 'D07-SMB-BASEMIAMIBEACH-PPSO@uscg.mil', now(), now());



delete from transportation_offices where name = 'PPPO Base San Pedro - USCG';