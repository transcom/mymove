CREATE TABLE IF NOT EXISTS gsr_appeals (
    id UUID PRIMARY KEY NOT NULL,
    evaluation_report_id UUID REFERENCES evaluation_reports(id) ON DELETE CASCADE,
    report_violation_id UUID REFERENCES report_violations(id) ON DELETE SET NULL,
    office_user_id UUID REFERENCES office_users(id) ON DELETE SET NULL NOT NULL,
    is_serious_incident_appeal BOOLEAN,
    appeal_status TEXT NOT NULL CHECK (appeal_status IN ('SUSTAINED', 'REJECTED')),
    remarks TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

COMMENT on TABLE gsr_appeals IS 'Stores appeals made from Government Surveillance Representative (GSR) users';
COMMENT on COLUMN gsr_appeals.evaluation_report_id IS 'Evaluation report that is associated with the appeal.';
COMMENT on COLUMN gsr_appeals.report_violation_id IS 'Report violation that is associated with the appeal.';
COMMENT on COLUMN gsr_appeals.office_user_id IS 'Office user that is leaving the appeal.';
COMMENT on COLUMN gsr_appeals.is_serious_incident_appeal IS 'Determines if the appeal is on a serious incident or not.';
COMMENT on COLUMN gsr_appeals.appeal_status IS 'Status of the appeal. Can be SUSTAINED or REJECTED.';
COMMENT on COLUMN gsr_appeals.remarks IS 'Remarks from GSR user when creating the appeal.';
COMMENT on COLUMN gsr_appeals.created_at IS 'Date that appeal was created.';
COMMENT on COLUMN gsr_appeals.updated_at IS 'Date that appeal was updated.';
COMMENT on COLUMN gsr_appeals.deleted_at IS 'Date that appeal was soft deleted.';

CREATE INDEX IF NOT EXISTS gsr_appeals_evaluation_report_id_idx ON gsr_appeals (evaluation_report_id);
CREATE INDEX IF NOT EXISTS gsr_appeals_report_violation_id_idx ON gsr_appeals (report_violation_id);
