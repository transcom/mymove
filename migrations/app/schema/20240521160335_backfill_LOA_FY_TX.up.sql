DO $$
BEGIN
    UPDATE lines_of_accounting
    SET loa_bg_fy_tx = CASE
                           WHEN extract(month FROM loa_bgn_dt) > 10 THEN extract(year FROM loa_bgn_dt)
                           ELSE extract(year FROM loa_bgn_dt) - 1
                       END
    WHERE loa_bg_fy_tx IS NULL;

    UPDATE lines_of_accounting
    SET loa_end_fy_tx = CASE
                            WHEN extract(month FROM loa_end_dt) > 10 THEN extract(year FROM loa_end_dt)
                            ELSE extract(year FROM loa_end_dt) - 1
                        END
    WHERE loa_end_fy_tx IS NULL;
END $$;
