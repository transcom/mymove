-- New Columns
ALTER TABLE weight_tickets
ADD COLUMN adjusted_net_weight int,
ADD COLUMN net_weight_remarks varchar;

-- Column Comments
COMMENT ON COLUMN weight_tickets.adjusted_net_weight IS 'Stores the net weight of the vehicle';
COMMENT ON COLUMN weight_tickets.net_weight_remarks IS 'Stores remarks explaining any edits made to the net weight';
