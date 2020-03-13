-- Change nullability
ALTER TABLE customers
    ALTER COLUMN user_id SET NOT NULL,
    ALTER COLUMN dod_id DROP NOT NULL;

-- Fix any nullable columns to store null instead of empty string
UPDATE customers SET dod_id = NULL WHERE dod_id = '';
UPDATE customers SET first_name = NULL WHERE first_name = '';
UPDATE customers SET last_name = NULL WHERE last_name = '';
UPDATE customers SET agency = NULL WHERE agency = '';
