ALTER TABLE entitlements
    ADD COLUMN pro_gear_weight int4;
COMMENT ON COLUMN entitlements.pro_gear_weight IS 'This is equipment a member needs for the performance of official duties at the next or a later destination. Members are given a weight allowance for progear that is separate from their normal weight allowance.';
ALTER TABLE entitlements
    ADD COLUMN spouse_pro_gear_weight int4;
COMMENT ON COLUMN entitlements.spouse_pro_gear_weight IS 'This is equipment a member''s spouse needs for the performance of official duties at the next or a later destination. Members are given a weight allowance for progear that is separate from their normal weight allowance.';
