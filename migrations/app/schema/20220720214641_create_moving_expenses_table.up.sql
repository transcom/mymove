CREATE TYPE moving_expense_type AS enum (
	'CONTRACTED_EXPENSE',
	'OIL',
	'PACKING_MATERIALS',
	'RENTAL_EQUIPMENT',
	'STORAGE',
	'TOLLS',
	'WEIGHING_FEES',
	'OTHER'
	);

CREATE TYPE ppm_document_status AS enum (
	'APPROVED',
	'EXCLUDED',
	'REJECTED'
	);

CREATE TABLE moving_expenses
(
	id uuid PRIMARY KEY,
	ppm_shipment_id uuid NOT NULL
		CONSTRAINT moving_expenses_ppm_shipments_id_fkey
			REFERENCES ppm_shipments,
	document_id uuid NOT NULL
		CONSTRAINT moving_expenses_document_id_fkey
			REFERENCES documents,
	moving_expense_type moving_expense_type,
	description varchar,
	paid_with_gtcc bool,
	amount int,
	missing_receipt bool,
	status ppm_document_status,
	sit_start_date date,
	sit_end_date date,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamptz
);

CREATE INDEX moving_expenses_ppm_shipment_id_idx ON moving_expenses USING hash (ppm_shipment_id);
CREATE INDEX moving_expenses_deleted_at_idx ON moving_expenses USING btree (deleted_at);

COMMENT on TABLE moving_expenses IS 'Stores weight ticket docs associated with a trip for a PPM shipment.';
COMMENT on COLUMN moving_expenses.ppm_shipment_id IS 'The ID of the PPM shipment that this set of weight tickets is for.';
COMMENT on COLUMN moving_expenses.document_id IS 'The ID of the document that is associated with the user uploads containing the moving expense receipt.';
COMMENT on COLUMN moving_expenses.moving_expense_type IS 'Identifies the type of expense this is.';
COMMENT on COLUMN moving_expenses.description IS 'Stores a description of the expense.';
COMMENT on COLUMN moving_expenses.paid_with_gtcc IS 'Indicates if the customer paid using a Government Travel Charge Card (GTCC).';
COMMENT on COLUMN moving_expenses.amount IS 'Stores the cost of the expense.';
COMMENT on COLUMN moving_expenses.missing_receipt IS 'Indicates if the customer is missing the receipt for their expense.';
COMMENT on COLUMN moving_expenses.status IS 'Status of the expense, e.g. APPROVED.';
COMMENT on COLUMN moving_expenses.sit_start_date IS 'If this is a STORAGE expense, this indicates the date storage began.';
COMMENT on COLUMN moving_expenses.sit_end_date IS 'If this is a STORAGE expense, this indicates the date storage ended.';
