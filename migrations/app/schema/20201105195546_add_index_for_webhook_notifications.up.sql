-- Recreate index
DROP INDEX IF EXISTS webhook_notifications_unsent;
CREATE INDEX webhook_notifications_unsent ON webhook_notifications(created_at)
    WHERE status != 'SENT' AND status != 'SKIPPED';
