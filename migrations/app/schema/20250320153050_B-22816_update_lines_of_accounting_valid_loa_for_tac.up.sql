--B-22816   Maria Traskowsky    Update lines_of_accounting.valid_loa_for_tac column

-- Set temp timeout due to large file modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

UPDATE lines_of_accounting
SET valid_loa_for_tac =
    CASE
        WHEN loa_dpt_id IS NOT NULL
         AND loa_baf_id IS NOT NULL
         AND loa_trsy_sfx_tx IS NOT NULL
         AND loa_obj_cls_id IS NOT NULL
         AND loa_doc_id IS NOT NULL
         AND loa_bgn_dt IS NOT NULL
         AND loa_end_dt IS NOT NULL
         AND loa_uic IS NOT NULL
        THEN TRUE
        ELSE FALSE
    END;