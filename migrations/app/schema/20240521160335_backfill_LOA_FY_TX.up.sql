DO $$
BEGIN
    UPDATE lines_of_accounting
    SET loa_bg_fy_tx = CASE
                           WHEN extract(month FROM loa_bgn_dt) > 9 THEN extract(year FROM loa_bgn_dt) +1
                           ELSE extract(year FROM loa_bgn_dt)
                       END
    WHERE loa_bg_fy_tx IS NULL;

    UPDATE lines_of_accounting
    SET loa_end_fy_tx = CASE
                            WHEN extract(month FROM loa_end_dt) > 9 THEN extract(year FROM loa_end_dt) +1
                            ELSE extract(year FROM loa_end_dt)
                        END
    WHERE loa_end_fy_tx IS NULL;
END $$;
