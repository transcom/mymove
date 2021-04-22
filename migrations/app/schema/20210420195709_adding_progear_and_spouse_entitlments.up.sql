ALTER TABLE entitlements
    ADD COLUMN pro_gear_weight integer DEFAULT 0 NOT NULL CHECK (pro_gear_weight >= 0 AND pro_gear_weight <= 2000);
COMMENT ON COLUMN entitlements.pro_gear_weight IS 'This is equipment a member needs for the performance of official duties at the next or a later destination. Members are given a weight allowance for progear that is separate from their normal weight allowance.';
ALTER TABLE entitlements
    ADD COLUMN pro_gear_weight_spouse integer DEFAULT 0 NOT NULL CHECK (pro_gear_weight_spouse >= 0 AND pro_gear_weight_spouse <= 500);
COMMENT ON COLUMN entitlements.pro_gear_weight_spouse IS 'This is equipment a member''s spouse needs for the performance of official duties at the next or a later destination. Members are given a weight allowance for progear that is separate from their normal weight allowance.';
