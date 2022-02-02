
CREATE TYPE ppm_shipment_status AS enum (
	'DRAFT',
	'SUBMITTED',
	'APPROVED',
    'PAYMENT_REQUESTED',
    'COMPLETED',
    'CANCELED'
	);

CREATE TABLE ppm_shipment
(
	id uuid PRIMARY KEY NOT NULL,
	shipment_id    uuid NOT NULL
		CONSTRAINT ppm_shipment_mto_shipment_id_fkey
		REFERENCES mto_shipments,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	status ppm_shipment_status NOT NULL,
	expected_departure_date date,
	actual_move_date date,
	submitted_at timestamptz,
	reviewed_at timestamptz,
	approved_at timestamptz,
	pickup_postal_code varchar,
	secondary_pickup_postal_code varchar,
	destination_postal_code varchar,
	secondary_destination_postal_code varchar,
	sit_expected bool,
	estimated_weight int,
	net_weight int,
	has_pro_gear bool,
	pro_gear_weight int,
	spouse_pro_gear_weight int,
	estimated_incentive int,
	advance_requested bool,
	advance_id uuid
		CONSTRAINT ppm_shipment_reimbursements_id_fkey
		REFERENCES reimbursements,
	advance_worksheet_id uuid
		CONSTRAINT ppm_shipment_documents_id_fkey
		REFERENCES documents
);

COMMENT on TABLE ppm_shipment IS 'Stores all PPM shipments, and their details.';
COMMENT on COLUMN ppm_shipment.shipment_id IS 'MTO shipment ID associated with this PPM shipment.';
COMMENT on COLUMN ppm_shipment.status IS 'Status of the PPM shipment.';
COMMENT on COLUMN ppm_shipment.expected_departure_date IS 'Expected date this PPM shipment begins.';
COMMENT on COLUMN ppm_shipment.actual_move_date IS 'Actual date of the move associated with this PPM shipment.';
COMMENT on COLUMN ppm_shipment.submitted_at IS 'Date that PPM shipment information was submitted.';
COMMENT on COLUMN ppm_shipment.reviewed_at IS 'Date that PPM shipment information was reviewed.';
COMMENT on COLUMN ppm_shipment.approved_at IS 'Date that PPM shipment information was approved.';
COMMENT on COLUMN ppm_shipment.pickup_postal_code IS 'Postal code where PPM begins.';
COMMENT on COLUMN ppm_shipment.secondary_pickup_postal_code IS 'Secondary postal code where PPM shipment is to be picked up.';
COMMENT on COLUMN ppm_shipment.destination_postal_code IS 'Destination postal code for PPM shipment.';
COMMENT on COLUMN ppm_shipment.secondary_destination_postal_code IS 'Secondary destination postal code for PPM shipment.';
COMMENT on COLUMN ppm_shipment.sit_expected IS 'Indicate if SIT is expected for PPM shipment.';
COMMENT on COLUMN ppm_shipment.estimated_weight IS 'Estimated weight of PPM shipment.';
COMMENT on COLUMN ppm_shipment.net_weight IS 'Net weight of PPM shipment.';
COMMENT on COLUMN ppm_shipment.has_pro_gear IS 'Indicate if PPM shipment has pro gear.';
COMMENT on COLUMN ppm_shipment.pro_gear_weight IS 'Indicate weight of PPM shipment pro gear.';
COMMENT on COLUMN ppm_shipment.spouse_pro_gear_weight IS 'Indicate weight of PPM shipment spouse pro gear.';
COMMENT on COLUMN ppm_shipment.estimated_incentive IS 'Estimated incentive associated with PPM shipment.';
COMMENT on COLUMN ppm_shipment.advance_requested IS 'Indicate if advance has been requested for PPM shipment.';
COMMENT on COLUMN ppm_shipment.advance_id IS 'Reimbursement ID for PPM shipment.';
COMMENT on COLUMN ppm_shipment.advance_worksheet_id IS 'Document ID for PPM shipment.';
