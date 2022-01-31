CREATE TABLE archived_personally_procured_move_data (
	id SERIAL PRIMARY KEY,
    move_id uuid REFERENCES moves,
	FOREIGN KEY (move_id) REFERENCES moves (id) ON DELETE CASCADE,
    weight_estimate INT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL,
    pickup_postal_code varchar(255),
    additional_pickup_postal_code varchar(255),
    destination_postal_code varchar(255),
    days_in_storage INT,
    status varchar(255),
    has_additional_postal_code bool,
    has_sit bool,
    has_requested_advance bool,
    advance_id uuid REFERENCES reimbursements,
    estimated_storage_reimbursement varchar(255),
    mileage INT,
    planned_sit_max INT,
    sit_max INT,
    incentive_estimate_min INT,
    incentive_estimate_max INT,
    advance_worksheet_id uuid REFERENCES documents,
    net_weight INT,
    original_move_date timestamp with time zone,
    actual_move_date timestamp with time zone,
    total_sit_cost INT,
    submit_date timestamp with time zone,
    approve_date timestamp with time zone,
    reviewed_date timestamp with time zone,
    has_pro_gear text,
    has_pro_gear_over_thousand text,
    move_document_id uuid REFERENCES move_documents,
    move_document_move_id uuid REFERENCES moves,
    document_id uuid REFERENCES documents,
    move_document_type varchar(255),
    move_document_status varchar(255),
    notes text,
    move_document_created_at timestamp with time zone,
    move_document_updated_at timestamp with time zone,
    move_document_deleted_at timestamp with time zone,
    title varchar(255),
    move_document_ppm_id uuid REFERENCES personally_procured_moves,
    signed_certificate_id uuid REFERENCES signed_certifications,
    submitting_user_id uuid REFERENCES users,
    signed_certificate_move_id uuid REFERENCES moves,
    certification_text text,
    signature text,
    date timestamp,
	signed_certificate_created_at timestamp without time zone,
	signed_certificate_updated_at timestamp without time zone,
    personally_procured_move_id uuid REFERENCES moves,
    certification_type text,
    weight_ticket_set_document_id uuid REFERENCES weight_ticket_set_documents,
    weight_ticket_set_type text,
    vehicle_nickname text,
    weight_ticket_set_document_move_document_id uuid REFERENCES move_documents,
    empty_weight INT,
    empty_weight_ticket_mising bool,
    full_weight INT,
    full_weight_ticket_missing bool,
    weight_ticket_date timestamp,
    weight_ticket_set_document_created_at timestamp without time zone,
    weight_ticket_set_document_updated_at timestamp without time zone,
    weight_ticket_set_document_deleted_at timestamp without time zone,
    trailer_ownership_missing bool,
    vehicle_make text,
    vehicle_model text,
    moving_expenses_document_id uuid REFERENCES moving_expense_documents,
    moving_expense_type varchar(255),
    moving_expenses_document_created_at timestamp without time zone,
    moving_expenses_document_updated_at timestamp without time zone,
    moving_expenses_document_deleted_at timestamp without time zone,
    requested_amount_cents INT,
    payment_method varchar(255),
    receipt_missing bool,
    storage_start_date timestamp,
    storage_end_date timestamp
);

CREATE INDEX ON archived_personally_procured_move_data (move_id);
CREATE INDEX ON archived_personally_procured_move_data (advance_id);
CREATE INDEX ON archived_personally_procured_move_data (advance_worksheet_id);
CREATE INDEX ON archived_personally_procured_move_data (move_document_id);
CREATE INDEX ON archived_personally_procured_move_data (move_document_move_id);
CREATE INDEX ON archived_personally_procured_move_data (document_id);
CREATE INDEX ON archived_personally_procured_move_data (move_document_ppm_id);
CREATE INDEX ON archived_personally_procured_move_data (signed_certificate_id);
CREATE INDEX ON archived_personally_procured_move_data (submitting_user_id);
CREATE INDEX ON archived_personally_procured_move_data (signed_certificate_move_id);
CREATE INDEX ON archived_personally_procured_move_data (personally_procured_move_id);
CREATE INDEX ON archived_personally_procured_move_data (weight_ticket_set_document_id);
CREATE INDEX ON archived_personally_procured_move_data (weight_ticket_set_document_move_document_id);
CREATE INDEX ON archived_personally_procured_move_data (moving_expenses_document_id);

INSERT INTO archived_personally_procured_move_data (id, move_id, weight_estimate, created_at, updated_at, pickup_postal_code, additional_pickup_postal_code,
destination_postal_code, days_in_storage, has_additional_postal_code, advance_id, has_requested_advance, advance_worksheet_id,
estimated_storage_reimbursement, mileage, planned_sit_max, sit_max, status, incentive_estimate_max, incentive_estimate_min,
net_weight, original_move_date, actual_move_date, total_sit_cost, submit_date, reviewed_date, approve_date,
has_pro_gear, has_pro_gear_over_thousand)
SELECT id, move_id, weight_estimate, created_at, updated_at, pickup_postal_code, additional_pickup_postal_code,
destination_postal_code, days_in_storage, has_additional_postal_code, advance_id, has_requested_advance, advance_worksheet_id,
estimated_storage_reimbursement, mileage, planned_sit_max, sit_max, status, incentive_estimate_max, incentive_estimate_min,
net_weight, original_move_date, actual_move_date, total_sit_cost, submit_date, reviewed_date, approve_date,
has_pro_gear, has_pro_gear_over_thousand
FROM personally_procured_moves;

INSERT INTO archived_personally_procured_move_data (signed_certificate_id, submitting_user_id, signed_certificate_move_id,
certification_text, signature, date, signed_certificate_created_at, signed_certificate_updated_at, certification_type,
personally_procured_move_id)
SELECT id AS signed_certificate_id, submitting_user_id, move_id AS signed_certificate_move_id,
certification_text, signature, date, created_at AS signed_certificate_created_at, updated_at AS signed_certificate_updated_at, certification_type,
personally_procured_move_id
FROM signed_certifications;