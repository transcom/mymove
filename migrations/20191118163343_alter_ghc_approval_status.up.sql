CREATE TYPE ghc_approval_status_2 AS ENUM (
    'APPROVED',
    'SUBMITTED',
    'REJECTED'
    );

ALTER TABLE move_task_orders
	ALTER COLUMN status TYPE ghc_approval_status_2
		USING (status::text::ghc_approval_status_2);

DROP TYPE ghc_approval_status;

ALTER TYPE ghc_approval_status_2 RENAME to ghc_approval_status;