ALTER TABLE re_zip5_rate_areas
    ADD COLUMN contract_id uuid NOT NULL
        CONSTRAINT re_zip5_rate_areas_contract_id_fkey
            REFERENCES re_contracts,
    ADD CONSTRAINT re_zip3_rate_areas_unique_key UNIQUE (contract_id, zip5),
    DROP CONSTRAINT re_zip5_rate_areas_zip5_key;