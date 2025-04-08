SELECT cron.schedule(
        'flag_sent_to_gex_for_review',
        '0 * * * 1-5',
        $$SELECT flag_sent_to_gex_for_review() $$
    );
-- runs hourly, Monday through Friday