ALTER table signed_certifications
    ADD COLUMN ppm_id uuid CONSTRAINT ppm_id_fkey REFERENCES ppm_shipments (id);

CREATE INDEX on signed_certifications (ppm_id);

COMMENT ON COLUMN signed_certifications.ppm_id IS 'Customer address to receive their W2 tax form';
