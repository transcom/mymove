ALTER table ppm_shipments
    ADD COLUMN w2_address_id uuid CONSTRAINT ppm_shipments_address_id_fkey REFERENCES addresses (id),
    ADD COLUMN final_incentive int;

CREATE INDEX on ppm_shipments (w2_address_id);

COMMENT ON COLUMN ppm_shipments.w2_address_id IS 'Customer address on their W2 tax form';
COMMENT ON COLUMN ppm_shipments.final_incentive IS 'The final calculated incentive for the PPM shipment. This does not include SIT as it is a reimbursement.';
