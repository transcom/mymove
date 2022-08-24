ALTER table ppm_shipments
    ADD COLUMN w2_address_id uuid CONSTRAINT ppm_shipments_address_id_fkey REFERENCES addresses (id);

CREATE INDEX on ppm_shipments (w2_address_id);

COMMENT ON COLUMN ppm_shipments.w2_address_id IS 'Customer address to receive their W2 tax form';
