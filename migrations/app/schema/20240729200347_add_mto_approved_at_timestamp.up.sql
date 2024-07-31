ALTER TABLE moves
ADD COLUMN approved_at timestamp with time zone;
COMMENT ON COLUMN moves.approved_at IS 'Timestamp when the Move Task Order had its shipments submitted by the Task Ordering Officer and is thus considered approved.';