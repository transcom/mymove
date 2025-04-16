--B-22540  Alex Lusk  Create ppm_closeouts table to store dollar values for Advana
--B-22545  Alex Lusk  Remove unused storage expense columns

CREATE TABLE IF NOT EXISTS ppm_closeouts (
    id UUID PRIMARY KEY NOT NULL,
    ppm_shipment_id UUID REFERENCES ppm_shipments(id) ON DELETE CASCADE NOT NULL,
    max_advance integer,
	gtcc_paid_contracted_expense integer,
	member_paid_contracted_expense integer,
	gtcc_paid_packing_materials integer,
	member_paid_packing_materials integer,
	gtcc_paid_weighing_fee integer,
	member_paid_weighing_fee integer,
	gtcc_paid_rental_equipment integer,
	member_paid_rental_equipment integer,
	gtcc_paid_tolls integer,
	member_paid_tolls integer,
	gtcc_paid_oil integer,
	member_paid_oil integer,
	gtcc_paid_other integer,
	member_paid_other integer,
	total_gtcc_paid_expenses integer,
	total_member_paid_expenses integer,
	remaining_incentive integer,
	gtcc_paid_sit integer,
	member_paid_sit integer,
	gtcc_disbursement integer,
	member_disbursement integer,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

COMMENT on TABLE ppm_closeouts IS 'Stores dollar values paid out in PPM closeouts that are present on the shipment summary worksheet';

COMMENT on COLUMN ppm_closeouts.ppm_shipment_id IS 'PPM shipment that is associated with this closeout.';
COMMENT on COLUMN ppm_closeouts.max_advance IS 'Maximum value a customer can request as an advance. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_contracted_expense IS 'Amount paid for contracted expenses using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_contracted_expense IS 'Amount paid for contracted expenses by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_packing_materials IS 'Amount paid for packing materials using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_packing_materials IS 'Amount paid for packing materials by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_weighing_fee IS 'Amount paid for weighing fees using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_weighing_fee IS 'Amount paid for weighing fees by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_rental_equipment IS 'Amount paid for rental equipment using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_rental_equipment IS 'Amount paid for rental equipment by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_tolls IS 'Amount paid for tolls using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_tolls IS 'Amount paid for tolls by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_oil IS 'Amount paid for oil using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_oil IS 'Amount paid for oil by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_other IS 'Amount paid for other expenses using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_other IS 'Amount paid for other expenses by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.total_gtcc_paid_expenses IS 'Total amount paid using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.total_member_paid_expenses IS 'Total amount paid by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.remaining_incentive IS 'Final PPM incentive less the advance recieved. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_paid_sit IS 'Amount paid for SIT using the service member''s GTCC. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_paid_sit IS 'Amount paid for SIT by the service member. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.gtcc_disbursement IS 'Amount disbursed for GTCC expenses. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.member_disbursement IS 'Amount disbursed for service member paid expenses. Stored in cents.';
COMMENT on COLUMN ppm_closeouts.created_at IS 'Date that this closeout was created.';
COMMENT on COLUMN ppm_closeouts.updated_at IS 'Date that this closeout was updated.';

CREATE INDEX IF NOT EXISTS ppm_closeouts_ppm_shipment_id_idx ON ppm_closeouts (ppm_shipment_id);

ALTER TABLE ppm_closeouts
	DROP COLUMN IF EXISTS gtcc_paid_storage,
	DROP COLUMN IF EXISTS member_paid_storage;
