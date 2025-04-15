--change counseling office PPPO McChord Field - USA to PPPO JB Lewis-McChord (McChord) - USA
update transportation_offices set name = 'PPPO JB Lewis-McChord (McChord) - USA' where id = '95abaeaa-452f-4fe0-9264-960cd2a15ccd';

--remove counseling office PPPO DMO Mountain Warfare Training Center Bridgeport – USMC
update moves m
   set counseling_transportation_office_id = '311b5292-6a8c-4ed4-a7e1-374734118737'
  from orders o
 where m.counseling_transportation_office_id = 'fab58a38-ee1f-4adf-929a-2dd246fc5e67'
   and m.orders_id = o.id
   and o.origin_duty_location_id = '74651905-dd53-49f9-a196-6c3e9b43c734';

update moves m
   set counseling_transportation_office_id = '3210a533-19b8-4805-a564-7eb452afce10'
  from orders o
 where m.counseling_transportation_office_id = 'fab58a38-ee1f-4adf-929a-2dd246fc5e67'
   and m.orders_id = o.id
   and o.origin_duty_location_id = 'd9410393-3166-478e-a991-0c666998277f';

update duty_locations set transportation_office_id = null where id = '74651905-dd53-49f9-a196-6c3e9b43c734';
delete from transportation_offices where id = 'fab58a38-ee1f-4adf-929a-2dd246fc5e67';

--update counseling office name for Camp Lejeune from PPPO to PPSO
update transportation_offices set name = 'PPSO DMO Camp Lejeune - USMC' where id = '22894aa1-1c29-49d8-bd1b-2ce64448cc8d';

--update counseling office name for PPPO Miami - USA
update transportation_offices set name = 'PPPO Miami - USA' where id = '7f7cc97c-2f3c-4866-90fe-b335f5c8e042';

--update city names per Danny
update re_cities set city_name = 'JB LEWIS MCCHORD' where id = '768b1d81-f7a5-4352-921d-2fbcd5ff1f7f';
update addresses set city = 'JB LEWIS MCCHORD' where us_post_region_cities_id = '1616850b-e70f-4bd6-9bc9-43dee24cda2d';

update re_cities set city_name = 'BUCKLEY SFB' where id = '0776da0b-0687-45c8-b2ed-d1742b0043fd';
update addresses set city = 'BUCKLEY SFB' where us_post_region_cities_id = '5176b234-dbdb-4489-b1cd-be8623ad7865';

update addresses set city = 'HOLLOMAN AFB', us_post_region_cities_id = 'f3f3e2ce-501b-4832-ba35-d82ffe5add9a' where us_post_region_cities_id = '39723540-3e63-44e2-acac-04fb7ea44276';

update re_cities set city_name = 'LEMOORE NAS' where id = '397c4595-c57c-44e3-99f2-29375b4227c5';
update addresses set city = 'LEMOORE NAS' where us_post_region_cities_id = '52d5197b-5e03-440b-9483-81c4007cf951';

update re_cities set city_name = 'PETERSON SFB' where id = '8e404dbb-0d24-44b9-9096-c05290dc46cf';
update addresses set city = 'PETERSON SFB' where us_post_region_cities_id = '304c6c9c-6384-4329-9a00-26a5c3e515e5';

update re_cities set city_name = 'VANDENBERG SFB' where id = '589a5cec-0439-4ce4-8633-01ea0f177f55';
update addresses set city = 'VANDENBERG SFB' where us_post_region_cities_id = '729ddcb6-e65a-4d53-bb72-de4ef89d87c5';

--add counseling office PPPO DMO MCAGCC 29 Palms
update duty_locations set transportation_office_id = 'bd733387-6b6c-42ba-b2c3-76c20cc65006' where id = '1b60502a-f6be-479b-9a77-c29ea7e47043';