-- B-23372 Brooklyn Welsh - Add actual_gun_safe_weight column to mto_shipments to track gun safe weight across all tickets related to the shipment
ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS actual_gun_safe_weight int CHECK (actual_gun_safe_weight >= 0) NULL;