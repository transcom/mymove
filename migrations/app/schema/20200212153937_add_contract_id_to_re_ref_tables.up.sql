ALTER TABLE re_domestic_service_areas
    ADD COLUMN contract_id uuid NOT NULL
        CONSTRAINT re_domestic_service_areas_contract_id_fkey
            REFERENCES re_contracts,
    ADD CONSTRAINT re_domestic_service_areas_unique_key UNIQUE (contract_id, service_area),
    DROP CONSTRAINT re_domestic_service_areas_service_area_key;

ALTER TABLE re_rate_areas
    ADD COLUMN contract_id uuid NOT NULL
        CONSTRAINT re_rate_areas_contract_id_fkey
            REFERENCES re_contracts,
    ADD CONSTRAINT re_rate_areas_unique_key UNIQUE (contract_id, code),
    DROP CONSTRAINT re_rate_areas_code_key;

ALTER TABLE re_zip3s
    ADD COLUMN contract_id uuid NOT NULL
        CONSTRAINT re_zip3s_contract_id_fkey
            REFERENCES re_contracts,
    ADD CONSTRAINT re_zip3s_unique_key UNIQUE (contract_id, zip3),
    DROP CONSTRAINT re_zip3s_zip3_key;
