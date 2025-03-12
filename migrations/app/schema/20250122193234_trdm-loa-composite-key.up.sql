SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

-- Identify and order any existing duplicates (Duplicates will have only existed based on old existing migration files. Our TRDM data ingester does not allow for duplicates)
WITH ordered_loas AS (
  SELECT
    id,
    ROW_NUMBER() OVER (
      PARTITION BY loa_sys_id
      ORDER BY created_at DESC
    ) AS row_num
  FROM lines_of_accounting
)

-- Delete duplicates, we only want the most recent
DELETE FROM lines_of_accounting
WHERE id IN (
  SELECT id
  FROM ordered_loas
  WHERE row_num > 1
);

-- Set LOA unique constraint
ALTER TABLE lines_of_accounting
    ADD CONSTRAINT lines_of_accounting_loa_sys_id_unique
    UNIQUE (loa_sys_id);