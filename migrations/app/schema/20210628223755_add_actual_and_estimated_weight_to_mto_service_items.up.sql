-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN estimated_weight integer,
    ADD COLUMN actual_weight integer;

-- Column comments
COMMENT ON COLUMN mto_service_items.estimated_weight IS 'The guessed weight for a shuttling service item. Relevant for DOSHUT service items.';
COMMENT ON COLUMN mto_service_items.actual_weight IS 'The measured weight for a shuttling service item. Relevant for DOSHUT service items.';
