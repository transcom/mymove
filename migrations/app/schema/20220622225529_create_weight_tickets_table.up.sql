CREATE TABLE weight_tickets
(
	id uuid PRIMARY KEY,
	ppm_shipment_id uuid NOT NULL
		CONSTRAINT weight_tickets_ppm_shipments_id_fkey
		REFERENCES ppm_shipments,
	vehicle_description varchar,
	empty_weight int,
	has_empty_weight_ticket bool,
	empty_document_id uuid
	    CONSTRAINT weight_tickets_empty_document_id_fkey
	    REFERENCES documents,
	full_weight int,
	has_full_weight_ticket bool,
	full_document_id uuid
	    CONSTRAINT weight_tickets_full_document_id_fkey
		REFERENCES documents,
	owns_trailer bool,
	trailer_meets_criteria bool,
	proof_of_trailer_ownership_document_id uuid
		CONSTRAINT weight_tickets_proof_of_trailer_ownership_document_id_fkey
		REFERENCES documents,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamptz
);

CREATE INDEX weight_tickets_ppm_shipment_id_idx ON weight_tickets USING hash (ppm_shipment_id);
CREATE INDEX weight_tickets_deleted_at_idx ON weight_tickets USING btree (deleted_at);

COMMENT on TABLE weight_tickets IS 'Stores weight ticket docs associated with a trip for a PPM shipment.';
COMMENT on COLUMN weight_tickets.ppm_shipment_id IS 'The ID of the PPM shipment that this set of weight tickets is for.';
COMMENT on COLUMN weight_tickets.vehicle_description IS 'Stores a description of the vehicle used for the trip. E.g. make/model, type of truck/van, etc.';
COMMENT on COLUMN weight_tickets.empty_weight IS 'Stores the weight of the vehicle when empty.';
COMMENT on COLUMN weight_tickets.has_empty_weight_ticket IS 'Indicates if the customer has a weight ticket for the vehicle weight when empty.';
COMMENT on COLUMN weight_tickets.empty_document_id IS 'The ID of the document that is associated with the user uploads containing the empty vehicle weight.';
COMMENT on COLUMN weight_tickets.full_weight IS 'Stores the weight of the vehicle when full.';
COMMENT on COLUMN weight_tickets.has_full_weight_ticket IS 'Indicates if the customer has a weight ticket for the vehicle weight when full.';
COMMENT on COLUMN weight_tickets.empty_document_id IS 'The ID of the document that is associated with the user uploads containing the full vehicle weight.';
COMMENT on COLUMN weight_tickets.owns_trailer IS 'Indicates if the customer used a trailer they own for the move.';
COMMENT on COLUMN weight_tickets.trailer_meets_criteria IS 'Indicates if the trailer that the customer used meets all the criteria to be claimable.';
COMMENT on COLUMN weight_tickets.proof_of_trailer_ownership_document_id IS 'The ID of the document that is associated with the user uploads containing the proof of trailer ownership.';
