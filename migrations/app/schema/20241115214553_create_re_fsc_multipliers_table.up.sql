CREATE TABLE IF NOT EXISTS re_fsc_multipliers (
    id          uuid             NOT NULL,
    low_weight  int              NOT NULL,
    high_weight int              NOT NULL,
    multiplier  decimal          NOT NULL,
    created_at  timestamp        NOT NULL DEFAULT NOW(),
    updated_at  timestamp        NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE re_fsc_multipliers IS 'Stores data needed to needed to calculate FSC';
COMMENT ON COLUMN re_fsc_multipliers.low_weight IS 'The lowest weight permitted for a shipment';
COMMENT ON COLUMN re_fsc_multipliers.high_weight IS 'The highest weight permitted for a shipment';
COMMENT ON COLUMN re_fsc_multipliers.multiplier IS 'The decimal multiplier used to calculate the FSC';