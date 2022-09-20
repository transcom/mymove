CREATE TABLE IF NOT EXISTS report_violations (
	id uuid PRIMARY KEY,
	report_id uuid NOT NULL
		CONSTRAINT report_violations_report_id_fkey
		REFERENCES "evaluation_reports",
	violation_id uuid NOT NULL
		CONSTRAINT report_violations_violation_id_fkey
		REFERENCES "pws_violations"
);

CREATE INDEX report_violations_report_id_idx ON report_violations (report_id);

COMMENT ON TABLE report_violations IS 'Associates PWS Violations with QAE evaluation report.';
COMMENT ON COLUMN report_violations.report_id IS 'Report ID of the report that violations will be assiocated with.';
COMMENT ON COLUMN report_violations.violation_id IS 'Violation ID of the violation that will be assiocated to a report.';
