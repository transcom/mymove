
-- Add enum values for FAILING and SKIPPED
ALTER TYPE webhook_notifications_status
    ADD VALUE 'FAILING';
ALTER TYPE webhook_notifications_status
    ADD VALUE 'SKIPPED';

-- Add first_attempted_at column
ALTER TABLE webhook_notifications
	ADD COLUMN first_attempted_at timestamp without time zone;

