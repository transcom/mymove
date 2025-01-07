--update duty location for NAS Meridian, MS to use zip 39309
update duty_locations set name = 'NAS Meridian, MS 39309', address_id = '691551c2-71fe-4a15-871f-0c46dff98230' where id = '334fecaf-abeb-49ce-99b5-81d69c8beae5';

--remove 39302 duty location
delete from duty_locations where id = 'e55be32c-bf89-4927-8893-4454a26bfd55';
