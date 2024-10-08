DO $$
BEGIN
    UPDATE addresses
    SET is_oconus = CASE
                        WHEN (Select country from re_countries where id = country_id ) IN ('US', 'United States', 'USA') AND state NOT IN ('AK', 'HI') THEN false
                        ELSE true
                    END;
END $$;