--remove duty loc Johnston City, TN 37602
update orders set origin_duty_location_id = 'cd0c7325-15bb-45c7-a690-26c56c903ed7' where origin_duty_location_id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687';
update orders set new_duty_location_id = 'cd0c7325-15bb-45c7-a690-26c56c903ed7' where new_duty_location_id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687';

delete from duty_locations where id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687';

--remove duty loc Oceanside, CA 92052
update orders set origin_duty_location_id = 'a6993e7b-4600-44b9-b288-04ca011143f0' where origin_duty_location_id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de';
update orders set new_duty_location_id = 'a6993e7b-4600-44b9-b288-04ca011143f0' where new_duty_location_id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de';

delete from duty_locations where id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de';

