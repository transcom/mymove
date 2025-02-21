--B-22651   Maria Traskowsky    Add ub_weight_restriction column to entitlements table
ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS ub_weight_restriction int;