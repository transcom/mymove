
CREATE TYPE ppm_shipment_status AS enum (
	'SUBMITTED',
	'WAITING_ON_CUSTOMER',
	'NEEDS_ADVANCE_APPROVAL',
    'NEEDS_PAYMENT_APPROVAL',
    'PAYMENT_APPROVED'
	);

CREATE TABLE ppm_shipments
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

COMMENT on TABLE ppm_shipments IS 'Stores all PPM shipments, and their details.';
COMMENT on COLUMN ppm_shipments.shipment_id IS 'MTO shipment ID associated with this PPM shipment.';
COMMENT on COLUMN ppm_shipments.status IS 'Status of the PPM shipment.';
COMMENT on COLUMN ppm_shipments.expected_departure_date IS 'Expected date this PPM shipment begins.';
COMMENT on COLUMN ppm_shipments.actual_move_date IS 'Actual date of the move associated with this PPM shipment.';
COMMENT on COLUMN ppm_shipments.submitted_at IS 'Date that PPM shipment information was submitted.';
COMMENT on COLUMN ppm_shipments.reviewed_at IS 'Date that PPM shipment information was reviewed.';
COMMENT on COLUMN ppm_shipments.approved_at IS 'Date that PPM shipment information was approved.';
COMMENT on COLUMN ppm_shipments.pickup_postal_code IS 'Postal code where PPM begins.';
COMMENT on COLUMN ppm_shipments.secondary_pickup_postal_code IS 'Secondary postal code where PPM shipment is to be picked up.';
COMMENT on COLUMN ppm_shipments.destination_postal_code IS 'Destination postal code for PPM shipment.';
COMMENT on COLUMN ppm_shipments.secondary_destination_postal_code IS 'Secondary destination postal code for PPM shipment.';
COMMENT on COLUMN ppm_shipments.sit_expected IS 'Indicate if SIT is expected for PPM shipment.';
COMMENT on COLUMN ppm_shipments.estimated_weight IS 'Estimated weight of PPM shipment.';
COMMENT on COLUMN ppm_shipments.net_weight IS 'Net weight of PPM shipment.';
COMMENT on COLUMN ppm_shipments.has_pro_gear IS 'Indicate if PPM shipment has pro gear.';
COMMENT on COLUMN ppm_shipments.pro_gear_weight IS 'Indicate weight of PPM shipment pro gear.';
COMMENT on COLUMN ppm_shipments.spouse_pro_gear_weight IS 'Indicate weight of PPM shipment spouse pro gear.';
COMMENT on COLUMN ppm_shipments.estimated_incentive IS 'Estimated incentive associated with PPM shipment.';
COMMENT on COLUMN ppm_shipments.advance_requested IS 'Indicate if advance has been requested for PPM shipment.';
COMMENT on COLUMN ppm_shipments.advance_id IS 'Reimbursement ID for PPM shipment.';
COMMENT on COLUMN ppm_shipments.advance_worksheet_id IS 'Document ID for PPM shipment.';
