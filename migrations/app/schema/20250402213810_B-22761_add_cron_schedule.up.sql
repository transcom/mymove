SELECT cron.schedule(
        'flag_sent_to_gex_for_review',
        '0 5 * * *',
        $$SELECT flag_sent_to_gex_for_review() $$
    );
-- runs every day at 5:00 UTC