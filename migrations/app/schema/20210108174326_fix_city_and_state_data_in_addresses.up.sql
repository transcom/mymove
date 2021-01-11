UPDATE addresses SET street_address_1 = TRIM(street_address_1), updated_at = now() WHERE street_address_1 <> TRIM(street_address_1);
UPDATE addresses SET street_address_2 = TRIM(street_address_2), updated_at = now() WHERE street_address_2 IS NOT NULL and street_address_2 <> TRIM(street_address_2);
UPDATE addresses SET street_address_3 = TRIM(street_address_3), updated_at = now() WHERE street_address_3 IS NOT NULL and street_address_3 <> TRIM(street_address_3);
UPDATE addresses SET city = TRIM(city), updated_at = now() WHERE city <> TRIM(city);
UPDATE addresses SET state = TRIM(state), updated_at = now() WHERE state <> TRIM(state);
