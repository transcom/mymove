ALTER TABLE evaluation_reports 
  ADD COLUMN time_depart time,
  ADD COLUMN eval_start time,
  ADD COLUMN eval_end time;

COMMENT ON COLUMN evaluation_reports.time_depart IS 'Time departed for the evaluation, recorded in 24 hour format without timezone info';
COMMENT ON COLUMN evaluation_reports.eval_start IS 'Time evaluation started, recorded in 24 hour format without timezone info';
COMMENT ON COLUMN evaluation_reports.eval_end IS 'Time evaluation ended, recorded in 24 hour format without timezone info';
  