--B-23372 Add actual_gun_safe_weight column to mto_shipments
ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS actual_gun_safe_weight int CHECK (actual_gun_safe_weight >= 0) NULL;