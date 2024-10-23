-- Update is_oconus value on the addresses table based on the addresses country and the state
UPDATE addresses
SET is_oconus = CASE
                    WHEN country_id is null OR (Select country from re_countries where id = country_id ) = 'US' AND state NOT IN ('AK', 'HI') THEN false
                    ELSE true
                END;

-- Update is_oconus value on the addresses table to be not null
ALTER TABLE addresses
    ALTER COLUMN is_oconus SET NOT NULL;
