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
