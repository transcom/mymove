CREATE TABLE transportation_accounting_codes
(
    id uuid NOT NULL,
    tac VARCHAR(4) UNIQUE NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX transportation_accounting_codes_tac_idx ON transportation_accounting_codes USING btree(tac);

-- comments on columns
COMMENT ON COLUMN "transportation_accounting_codes"."tac" IS 'A 4-digit alphanumeric transportation accounting code used to look up long lines of accounting.  These values are sourced from TGET.';

