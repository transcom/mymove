-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN estimated_weight integer,
    ADD COLUMN actual_weight integer;

-- Column comments
COMMENT ON COLUMN mto_service_items.estimated_weight IS 'An estimate of how much weight from a shipment will be included in a shuttling (DDSHUT & DOSHUT) service item.';
COMMENT ON COLUMN mto_service_items.actual_weight IS 'Provided by the movers, based on weight tickets. Relevant for shuttling (DDSHUT & DOSHUT) service items.';
