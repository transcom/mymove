-- These type migrations are necessary because they do not align with TRDM data
-- It is going from int(4) to varchar (20) so there is no risk of the integer number being too large for varchar 20

-- Cast loa_sys_id from int to varchar
ALTER TABLE lines_of_accounting
ALTER COLUMN loa_sys_id TYPE VARCHAR(20) USING CAST(loa_sys_id AS VARCHAR(20));
