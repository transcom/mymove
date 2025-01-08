--update duty location for NAS Meridian, MS to use zip 39309
update duty_locations set name = 'NAS Meridian, MS 39309', address_id = '691551c2-71fe-4a15-871f-0c46dff98230' where id = '334fecaf-abeb-49ce-99b5-81d69c8beae5';

--remove 39302 duty location
delete from duty_locations where id = 'e55be32c-bf89-4927-8893-4454a26bfd55';

--update duty location for Minneapolis, MN 55460 to use 55467
update orders set new_duty_location_id = 'fc4d669f-594a-4784-9831-bf2eb9f8948b' where new_duty_location_id = '4c960096-1fbc-4b9d-b7d9-5979a3ba7344';

--remove 55460 duty location
delete from duty_locations where id = '4c960096-1fbc-4b9d-b7d9-5979a3ba7344';