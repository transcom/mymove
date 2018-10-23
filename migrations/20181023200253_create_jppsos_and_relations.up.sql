
-- NAVSUP FLC Puget Sound Whidbey
update transportation_offices
 set gbloc='JENQ' where id = '2726d99c-eeaa-4c54-b24f-692ba0a78e2b';

update transportation_offices
 set shipping_office_id= '2726d99c-eeaa-4c54-b24f-692ba0a78e2b' where gbloc='JENQ' and id <> '2726d99c-eeaa-4c54-b24f-692ba0a78e2b';


--AGFM
insert into addresses
    (id, street_address_1, city, state, postal_code, created_at, updated_at)
values
    ('414c9063-5af7-4709-b1f2-d99464f80866', '25 Chennault Street Building 1723', 'Hanscom AFB', 'MA', '01731', now(), now());
insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
values
    ('3132b512-1889-4776-a666-9c08a24afe20' , 'JPPSO Northeast Detachment 2', 'AGFM', '414c9063-5af7-4709-b1f2-d99464f80866', 42.4579955, -71.27431229999999, now(), now());

update transportation_offices
 set shipping_office_id= '3132b512-1889-4776-a666-9c08a24afe20' where gbloc='AGFM' and id <> '3132b512-1889-4776-a666-9c08a24afe20';

--HAFC
insert into addresses
    (id, street_address_1, city, state, postal_code, created_at, updated_at)
values
    ('46a4805f-bc93-4639-8224-84c4db946ab8', '2261 Hughes Ave, Suite 160', 'Lackland AFB', 'TX', '78236', now(), now());
insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
values
    ('c2c440ae-5394-4483-84fb-f872e32126bb' , 'JPPSO South Central', 'HAFC', '46a4805f-bc93-4639-8224-84c4db946ab8', 29.3877765, -98.6205193
, now(), now());
update transportation_offices
 set shipping_office_id= 'c2c440ae-5394-4483-84fb-f872e32126bb' where gbloc='HAFC' and id <> 'c2c440ae-5394-4483-84fb-f872e32126bb';

-- KKFA
insert into addresses
    (id, street_address_1, city, state, postal_code, created_at, updated_at)
values
    ('69e76812-1361-413f-abdd-bf14c05e1b8b', '121 S. Tejon St., Suite 800', 'Colorado Springs', 'CO', '80903', now(), now());
insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
values
    ('171b54fa-4c89-45d8-8111-a2d65818ff8c' , 'JPPSO North Central', 'KKFA', '69e76812-1361-413f-abdd-bf14c05e1b8b', 38.82066, -104.7051, now(), now());
update transportation_offices
 set shipping_office_id= '171b54fa-4c89-45d8-8111-a2d65818ff8c' where gbloc='KKFA' and id <> '171b54fa-4c89-45d8-8111-a2d65818ff8c';
