SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

-- Identify and order any existing duplicates (Duplicates will have only existed based on old existing migration files. Our TRDM data ingester does not allow for duplicates)
WITH ordered_tacs AS (
  SELECT
    id,
    ROW_NUMBER() OVER (
      PARTITION BY tac, tac_fy_txt, tac_fn_bl_mod_cd
      ORDER BY created_at DESC
    ) AS row_num
  FROM transportation_accounting_codes
)

-- Delete duplicates, we only want the most recent
DELETE FROM transportation_accounting_codes
WHERE id IN (
  SELECT id
  FROM ordered_tacs
  WHERE row_num > 1
);

-- Set TAC composite key, never to allow duplicates again
ALTER TABLE transportation_accounting_codes
    ADD CONSTRAINT transportation_accounting_codes_tac_fy_fbmc_unique
    UNIQUE (tac, tac_fy_txt, tac_fn_bl_mod_cd);