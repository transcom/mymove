-- These type migrations are necessary because they do not align with TRDM data
-- It is going from int(4) to varchar (20) so there is no risk of the integer number being too large for varchar 20
-- This is for all three column modifications

-- Cast tac_sys_id from int to varchar
ALTER TABLE transportation_accounting_codes
ALTER COLUMN tac_sys_id TYPE VARCHAR(20) USING CAST(tac_sys_id AS VARCHAR(20));

-- Cast loa_sys_id column from int to varchar
ALTER TABLE transportation_accounting_codes
ALTER COLUMN loa_sys_id TYPE VARCHAR(20) USING CAST(loa_sys_id AS VARCHAR(20));

-- Cast tac_fy_txt column from int to varchar
ALTER TABLE transportation_accounting_codes
ALTER COLUMN tac_fy_txt TYPE VARCHAR(20) USING CAST(tac_fy_txt AS VARCHAR(20));
