DO $$
BEGIN
    UPDATE addresses
    SET is_oconus = CASE
                        WHEN (Select country from re_countries where id = country_id ) = 'US' THEN false
                        ELSE true
                    END;
END $$;