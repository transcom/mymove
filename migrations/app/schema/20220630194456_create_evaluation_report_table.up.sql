CREATE TYPE evaluation_report_type AS enum (
	'DATA_REVIEW',
	'PHYSICAL',
	'VIRTUAL'
	);

CREATE TYPE evaluation_location_type AS enum (
	'ORIGIN',
	'DESTINATION',
	'OTHER'
	);

CREATE TABLE IF NOT EXISTS evaluation_reports
(
	id                        uuid PRIMARY KEY,
	office_user_id            uuid REFERENCES office_users NOT NULL,
	move_id                   uuid REFERENCES moves        NOT NULL,
	shipment_id               uuid REFERENCES mto_shipments,
	inspection_date           date,
	type                      evaluation_report_type,
	travel_time_minutes       int,
	location                  evaluation_location_type,
	location_description      text,
	observed_date             date,
	evaluation_length_minutes int,
	violations_observed       bool,
	remarks                   text,
	submitted_at              timestamp WITH TIME ZONE,
	created_at                timestamp WITH TIME ZONE     NOT NULL,
	updated_at                timestamp WITH TIME ZONE     NOT NULL
);

CREATE INDEX evaluation_reports_office_user_id_idx ON evaluation_reports (office_user_id);
CREATE INDEX evaluation_reports_move_id_idx ON evaluation_reports (move_id);
CREATE INDEX evaluation_reports_shipment_id_idx ON evaluation_reports (shipment_id);
CREATE INDEX evaluation_reports_submitted_at_idx ON evaluation_reports (submitted_at);

COMMENT ON TABLE evaluation_reports IS 'Contains QAE evaluation reports. There are two kinds of reports: shipment and counseling. You can tell them apart based on whether shipment_id is NULL.';
COMMENT ON COLUMN evaluation_reports.move_id IS 'Move that the report is associated with';
COMMENT ON COLUMN evaluation_reports.office_user_id IS 'The office_user who authored the evaluation report.';
COMMENT ON COLUMN evaluation_reports.shipment_id IS 'If present, indicates the shipment that this report is based on. NULL if this is not a shipment report.';
COMMENT ON COLUMN evaluation_reports.inspection_date IS 'date of inspection';
COMMENT ON COLUMN evaluation_reports.type IS 'Indicates the type of evaluation that is being described by this report. Either physical, virtual, or data review';
COMMENT ON COLUMN evaluation_reports.travel_time_minutes IS 'Amount of time that the evaluator spent travelling to and from the site to perform a physical evaluation';
COMMENT ON COLUMN evaluation_reports.location IS 'Indicates whether the inspection was performed at the origin or destination of the shipment. If OTHER is selected, location_description should contain a description of the alternative location';
COMMENT ON COLUMN evaluation_reports.location_description IS 'If the inspection was performed at a location other than the origin or destination of the shipment, this field contains a description of the location';
COMMENT ON COLUMN evaluation_reports.observed_date IS 'Date of pickup (if location=ORIGIN) or delivery (if location=DESTINATION), if it did not match the scheduled date. If pickup/delivery happened on the scheduled date, this will be left NULL';
COMMENT ON COLUMN evaluation_reports.evaluation_length_minutes IS 'Length of time spent performing the evaluation, in minutes';
COMMENT ON COLUMN evaluation_reports.violations_observed IS 'True if any PWS violations were observed during the inspection';
COMMENT ON COLUMN evaluation_reports.remarks IS 'Free text field for the evaluator''s notes about the inspection';
COMMENT ON COLUMN evaluation_reports.submitted_at IS 'Time when the report was submitted. If NULL, then the report is still a draft';







