ALTER TABLE ppm_shipments
	ADD COLUMN aoa_packet_id uuid
		CONSTRAINT ppm_shipments_aoa_packet_id_fkey
			REFERENCES documents,
	ADD COLUMN payment_packet_id uuid
		CONSTRAINT ppm_shipments_payment_packet_id_fkey
			REFERENCES documents;

COMMENT ON COLUMN ppm_shipments.aoa_packet_id IS 'The ID of the document that is associated with the upload containing the generated AOA packet for this PPM Shipment.';
COMMENT ON COLUMN ppm_shipments.payment_packet_id IS 'The ID of the document that is associated with the upload containing the generated payment packet for this PPM Shipment.';
