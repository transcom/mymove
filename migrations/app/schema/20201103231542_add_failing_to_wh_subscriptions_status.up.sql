-- Add enum values for FAILING and SKIPPED
ALTER TYPE webhook_subscriptions_status
    ADD VALUE 'FAILING';

-- Add severity column
ALTER TABLE webhook_subscriptions
	ADD COLUMN severity INT NOT NULL DEFAULT 0;

