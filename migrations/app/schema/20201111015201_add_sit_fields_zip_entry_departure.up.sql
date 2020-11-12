-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN sit_postal_code text,
    ADD COLUMN sit_entry_date date,
    ADD COLUMN sit_departure_date date;

-- Column comments
COMMENT ON COLUMN mto_service_items.sit_postal_code IS 'The postal code for the origin SIT facility where the Prime stores the shipment, used in pricing.';
COMMENT ON COLUMN mto_service_items.sit_entry_date IS 'The date when the Prime contractor places the shipment into a SIT facility. Relevant for DOFSIT and DDFSIT service items.';
COMMENT ON COLUMN mto_service_items.sit_departure_date IS 'The date when the Prime contractor removes the item from the SIT facility. Relevant for DOPSIT and DDDSIT service items.';
