-- Adds new columns to mto shipments table
ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS actual_pro_gear_weight INTEGER NULL,
ADD COLUMN IF NOT EXISTS actual_spouse_pro_gear_weight INTEGER NULL;

-- Comments on new columns
COMMENT ON COLUMN mto_shipments.actual_pro_gear_weight IS 'Indicates weight of MTO shipment pro gear.';
COMMENT ON COLUMN mto_shipments.actual_spouse_pro_gear_weight IS 'Indicates weight of MTO shipment spouse pro gear.';