-- only schedule cron job if pg_cron extension exists
DO $do$ BEGIN IF EXISTS (
    SELECT 1
    FROM pg_extension
    WHERE extname = 'pg_cron'
) THEN PERFORM cron.schedule(
    'flag_sent_to_gex_for_review',
    '0 * * * 1-5',
    -- runs hourly, Monday through Friday
    $$SELECT flag_sent_to_gex_for_review() $$
);
END IF;
END $do$;