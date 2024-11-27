-- enable tracking of the timestamp a move with ub shipments are at risk of excess, as well as
-- the timestamp of when an office user dismisses the notification

ALTER TABLE moves ADD COLUMN IF NOT EXISTS excess_unaccompanied_baggage_weight_qualified_at timestamp with time zone;
ALTER TABLE moves ADD COLUMN IF NOT EXISTS excess_unaccompanied_baggage_weight_acknowledged_at timestamp with time zone;

COMMENT ON COLUMN mto_shipments.excess_unaccompanied_baggage_weight_qualified_at IS "The date and time the sum of all the move's unaccompanied baggage shipment weights met or exceeded the excess unaccompanied baggage weight qualification threshold.";
COMMENT ON COLUMN mto_shipments.excess_unaccompanied_baggage_weight_acknowledged_at IS 'The date and time the TOO dismissed the risk of excess unaccompanied baggage weight alert.';
