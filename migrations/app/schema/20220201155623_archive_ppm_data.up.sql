CREATE TABLE archived_personally_procured_moves (
    id uuid PRIMARY KEY,
    size varchar(255),
    move_id uuid REFERENCES moves,
    weight_estimate INT,
	created_at timestamp without time zone,
	updated_at timestamp without time zone,
    pickup_postal_code varchar(255),
    additional_pickup_postal_code varchar(255),
    destination_postal_code varchar(255),
    days_in_storage INT,
    has_additional_postal_code bool,
    advance_id uuid REFERENCES reimbursements,
    has_requested_advance bool,
    advance_worksheet_id uuid REFERENCES documents,
    estimated_storage_reimbursement varchar(255),
    mileage INT,
    planned_sit_max INT,
    sit_max INT,
    status varchar(255),
    incentive_estimate_min INT,
    incentive_estimate_max INT,
    net_weight INT,
    original_move_date timestamp with time zone,
    actual_move_date timestamp with time zone,
    total_sit_cost INT,
    has_sit bool,
    submit_date timestamp with time zone,
    reviewed_date timestamp with time zone,
    approve_date timestamp with time zone,
    has_pro_gear text,
    has_pro_gear_over_thousand text
);

CREATE INDEX ON archived_personally_procured_moves (move_id);
CREATE INDEX ON archived_personally_procured_moves (advance_id);
CREATE INDEX ON archived_personally_procured_moves (advance_worksheet_id);

CREATE TABLE archived_move_documents(
    id uuid PRIMARY KEY,
    move_id uuid REFERENCES moves,
    document_id uuid REFERENCES documents,
    type varchar(255),
    status varchar(255),
    move_document_type varchar(255),
    notes text,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone NOT NULL,
    title varchar(255),
    personally_procured_move_id uuid REFERENCES archived_personally_procured_moves
);

CREATE INDEX ON archived_move_documents (id);
CREATE INDEX ON archived_move_documents (move_id);
CREATE INDEX ON archived_move_documents (document_id);
CREATE INDEX ON archived_move_documents (personally_procured_move_id);

CREATE TABLE archived_signed_certifications(
    id uuid PRIMARY KEY,
    submitting_user_id uuid REFERENCES users,
    move_id uuid REFERENCES moves,
    certification_text text,
    signature text,
    date timestamp,
	created_at timestamp without time zone,
	updated_at timestamp without time zone,
    personally_procured_move_id uuid REFERENCES archived_personally_procured_moves,
    certification_type text
);

CREATE INDEX ON archived_signed_certifications (id);
CREATE INDEX ON archived_signed_certifications (submitting_user_id);
CREATE INDEX ON archived_signed_certifications (move_id);
CREATE INDEX ON archived_signed_certifications (personally_procured_move_id);

CREATE TABLE archived_weight_ticket_set_documents(
    id uuid PRIMARY KEY,
    weight_ticket_set_type text,
    vehicle_nickname text,
    move_document_id uuid REFERENCES archived_move_documents,
    empty_weight INT,
    empty_weight_ticket_missing bool,
    full_weight INT,
    full_weight_ticket_missing bool,
    weight_ticket_date timestamp,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    trailer_ownership_missing bool,
    vehicle_make text,
    vehicle_model text
);

CREATE INDEX ON archived_weight_ticket_set_documents (id);
CREATE INDEX ON archived_weight_ticket_set_documents (move_document_id);

CREATE TABLE archived_moving_expense_documents(
    id uuid PRIMARY KEY,
    move_document_id uuid REFERENCES archived_move_documents,
    moving_expense_type varchar(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    requested_amount_cents INT,
    payment_method varchar(255),
    receipt_missing bool,
    storage_start_date timestamp,
    storage_end_date timestamp
);

CREATE INDEX ON archived_moving_expense_documents (id);
CREATE INDEX ON archived_moving_expense_documents (move_document_id);

INSERT INTO archived_personally_procured_moves (
id, size, move_id, weight_estimate,
created_at, updated_at,
pickup_postal_code, additional_pickup_postal_code,
destination_postal_code, days_in_storage, has_additional_postal_code,
advance_id, has_requested_advance, advance_worksheet_id,
estimated_storage_reimbursement, mileage, planned_sit_max, sit_max, status,
incentive_estimate_max, incentive_estimate_min,
net_weight, original_move_date, actual_move_date,
total_sit_cost, has_sit, submit_date, reviewed_date, approve_date,
has_pro_gear, has_pro_gear_over_thousand)
SELECT id, size, move_id, weight_estimate,
created_at, updated_at,
pickup_postal_code,
additional_pickup_postal_code,
destination_postal_code, days_in_storage,
has_additional_postal_code, advance_id,
has_requested_advance, advance_worksheet_id,
estimated_storage_reimbursement, mileage,
planned_sit_max, sit_max,
status,
incentive_estimate_max,
incentive_estimate_min,
net_weight, original_move_date,
actual_move_date, total_sit_cost,
has_sit, submit_date,
reviewed_date, approve_date,
has_pro_gear, has_pro_gear_over_thousand
FROM personally_procured_moves;

INSERT INTO archived_signed_certifications(
id, submitting_user_id, move_id,
certification_text, signature, date, created_at,
updated_at, certification_type, personally_procured_move_id)
SELECT id, submitting_user_id, move_id,
certification_text, signature, date,
created_at, updated_at, certification_type,
personally_procured_move_id
FROM signed_certifications;

INSERT INTO archived_move_documents(
id, move_id, document_id,
move_document_type, status, notes,
updated_at, created_at,
title, personally_procured_move_id, deleted_at)
SELECT id, move_id,
document_id, move_document_type, status, notes,
updated_at, created_at,
title, personally_procured_move_id, deleted_at
FROM move_documents;

INSERT INTO archived_weight_ticket_set_documents(
id, weight_ticket_set_type, vehicle_nickname,
move_document_id, empty_weight,
empty_weight_ticket_missing, full_weight, full_weight_ticket_missing, weight_ticket_date,
created_at, updated_at, trailer_ownership_missing,
deleted_at, vehicle_make, vehicle_model
)
SELECT id, weight_ticket_set_type,
vehicle_nickname, move_document_id,
empty_weight,
empty_weight_ticket_missing, full_weight,
full_weight_ticket_missing, weight_ticket_date,
created_at, updated_at, trailer_ownership_missing,
deleted_at, vehicle_make, vehicle_model
FROM weight_ticket_set_documents;

INSERT INTO archived_moving_expense_documents(
id, move_document_id, moving_expense_type,
created_at, updated_at, requested_amount_cents,
payment_method, receipt_missing, storage_start_date,
storage_end_date, deleted_at
)
SELECT id, move_document_id,
moving_expense_type,
created_at, updated_at,
requested_amount_cents,
payment_method, receipt_missing,
storage_start_date, storage_end_date,
deleted_at
from moving_expense_documents;
