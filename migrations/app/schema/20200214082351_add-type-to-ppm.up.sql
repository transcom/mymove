CREATE TYPE ppm_type AS ENUM (
    'FULL',
    'PARTIAL'
    );

ALTER TABLE personally_procured_moves
    ADD COLUMN type ppm_type;