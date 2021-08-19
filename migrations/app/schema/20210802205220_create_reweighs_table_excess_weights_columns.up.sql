CREATE TYPE reweigh_requester AS enum (
    'CUSTOMER',
    'PRIME',
    'SYSTEM',
    'TOO'
);

CREATE TABLE reweighs
(
    id uuid NOT NULL primary key,
    shipment_id uuid NOT NULL
        CONSTRAINT reweighs_shipment_id_fkey
            REFERENCES mto_shipments,
    requested_at timestamp with time zone NOT NULL,
    requested_by reweigh_requester NOT NULL,
    weight integer,
    verification_reason text,
    verification_provided_at timestamp with time zone,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE INDEX reweighs_shipment_id_idx ON reweighs (shipment_id);

COMMENT ON TABLE reweighs IS 'A reweigh represents a request from different users or the system for a shipment to be reweighed by the movers';
COMMENT ON COLUMN reweighs.shipment_id IS 'A foreign key that points to the mto_shipments table for which shipment is being reweighed. There should only be one reweigh request per shipment';
COMMENT ON COLUMN reweighs.requested_at IS 'The date and time when the reweigh request was initiated';
COMMENT ON COLUMN reweighs.requested_by IS 'The type of user who requested the reweigh, including automated requests determined by the milmove system';
COMMENT ON COLUMN reweighs.weight IS 'The reweighed weight in pounds (lbs) of the shipment submitted by the movers';
COMMENT ON COLUMN reweighs.verification_reason IS 'If a reweigh was requested but was not able to be performed the movers can provide an explanation';
COMMENT ON COLUMN reweighs.verification_provided_at IS 'The date and time when the verification_reason value was added';

ALTER TABLE moves ADD COLUMN excess_weight_qualified_at timestamp with time zone;

COMMENT ON COLUMN moves.excess_weight_qualified_at IS 'The date and time the sum of all the move''s shipments met the excess weight qualification threshold';

ALTER TABLE moves
    ADD COLUMN excess_weight_upload_id uuid
        CONSTRAINT moves_excess_weight_upload_id_fkey REFERENCES uploads;

COMMENT ON COLUMN moves.excess_weight_upload_id IS 'An uploaded document by the movers proving that the customer has been counseled about excess weight';

ALTER TABLE mto_shipments
    ADD COLUMN billable_weight_cap integer;

COMMENT ON COLUMN mto_shipments.billable_weight_cap IS 'The billable weight cap that the TIO can set per shipment that affects pricing';

ALTER TABLE mto_shipments
    ADD COLUMN billable_weight_justification text;

COMMENT ON COLUMN mto_shipments.billable_weight_justification IS 'The reasoning for why the TIO has set the billable_weight_cap to the chosen value';
