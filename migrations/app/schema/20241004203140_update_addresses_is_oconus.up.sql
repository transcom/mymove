DO $$
BEGIN
    UPDATE addresses
    SET is_oconus = CASE
                        WHEN country IN ('US', 'United States') THEN false
                        ELSE true
                    END;
END $$;