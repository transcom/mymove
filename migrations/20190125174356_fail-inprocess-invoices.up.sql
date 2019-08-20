-- With the first iteration of invoice submission code, submission failures left the
-- invoice in an IN_PROCESS state and there was no way to clean it up. That code has been fixed,
-- but this migration cleans up any remaining stragglers in the database.
-- Theoretically this could run against invoices that are ACTUALLY in progress, so we don't run this
-- against invoices that were created within the last 30 seconds.
UPDATE invoices
    SET status = 'SUBMISSION_FAILURE'
    WHERE status = 'IN_PROCESS'
    AND created_at < now() - INTERVAL '30 seconds';
