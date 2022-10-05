ALTER TABLE evaluation_reports
	ADD COLUMN serious_incident bool,
	ADD COLUMN serious_incident_desc text,
	ADD COLUMN observed_claims_response_date date,
	ADD COLUMN observed_pickup_date date,
	ADD COLUMN observed_pickup_spread_start_date date,
	ADD COLUMN observed_pickup_spread_end_date date;

COMMENT ON COLUMN evaluation_reports.serious_incident IS 'Indicates is a serious incident was found';
COMMENT ON COLUMN evaluation_reports.serious_incident_desc IS 'Text field for the description of the serious incident';
COMMENT ON COLUMN evaluation_reports.observed_claims_response_date IS 'Date of observed claims response ';
COMMENT ON COLUMN evaluation_reports.observed_pickup_date IS 'Date of observed pickup';
COMMENT ON COLUMN evaluation_reports.observed_pickup_spread_start_date IS 'Start date of observed pickup spread';
COMMENT ON COLUMN evaluation_reports.observed_pickup_spread_end_date IS 'End date of observed pickup spread';
