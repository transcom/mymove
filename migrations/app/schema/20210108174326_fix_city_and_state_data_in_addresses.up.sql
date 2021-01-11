UPDATE addresses SET city = TRIM(city), updated_at = now() WHERE city <> TRIM(city);
UPDATE addresses SET state = TRIM(state), updated_at = now() WHERE state <> TRIM(state);
