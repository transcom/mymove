
ALTER TABLE entitlements
ADD COLUMN accompanied_tour BOOLEAN NULL,
    ADD COLUMN IF NOT EXISTS dependents_under_twelve INTEGER NULL,
    ADD COLUMN IF NOT EXISTS dependents_twelve_and_over INTEGER NULL,
    ADD COLUMN IF NOT EXISTS ub_allowance INTEGER NULL;

COMMENT ON COLUMN entitlements.accompanied_tour IS 'Indicates if the move entitlement allows dependents to travel to the new Permanent Duty Station (PDS). This is only present on OCONUS moves.';
COMMENT ON COLUMN entitlements.dependents_under_twelve IS 'Indicates the number of dependents under the age of twelve for a move. This is only present on OCONUS moves.';
COMMENT ON COLUMN entitlements.dependents_twelve_and_over IS 'Indicates the number of dependents of the age twelve or older for a move. This is only present on OCONUS moves.';
COMMENT ON COLUMN entitlements.ub_allowance IS 'The amount of weight in pounds that the move is entitled for shipment types of Unaccompanied Baggage.';