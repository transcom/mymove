-------------------------------------------
-- 'duty_locations' dependents

-- NOTE: 'orders.origin/new_duty_location_id are intentionally excluded
-- to prevent accidental unintentional deletions.
-- Therefore, orders must be deleted prior to a duty location being deleted.
-------------------------------------------
-- Add ON DELETE CASCADE constraint for duty_location_names.duty_location_id
ALTER TABLE duty_location_names
DROP CONSTRAINT duty_location_names_duty_location_id_fkey;

ALTER TABLE duty_location_names
ADD CONSTRAINT duty_location_names_duty_location_id_fkey
FOREIGN KEY (duty_location_id)
REFERENCES duty_locations(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'orders' dependents
-------------------------------------------
-- Add ON DELETE CASCADE constraint for moves.order_id
ALTER TABLE moves
DROP CONSTRAINT moves_orders_id_fk;

ALTER TABLE moves
ADD CONSTRAINT moves_orders_id_fk
FOREIGN KEY (orders_id)
REFERENCES orders(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'moves' dependents
-------------------------------------------
-- Add ON DELETE CASCADE constraint for customer_support_remarks.move_id
ALTER TABLE customer_support_remarks
DROP CONSTRAINT fk_moves;

ALTER TABLE customer_support_remarks
ADD CONSTRAINT fk_moves
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for evaluation_reports.move_id
ALTER TABLE evaluation_reports
DROP CONSTRAINT evaluation_reports_move_id_fkey;

ALTER TABLE evaluation_reports
ADD CONSTRAINT evaluation_reports_move_id_fkey
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for personally_procured_moves.move_id
ALTER TABLE personally_procured_moves
DROP CONSTRAINT personally_procured_moves_move_id_fkey;

ALTER TABLE personally_procured_moves
ADD CONSTRAINT personally_procured_moves_move_id_fkey
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for signed_certifications.move_id
ALTER TABLE signed_certifications
DROP CONSTRAINT signed_certifications_moves_id_fk;

ALTER TABLE signed_certifications
ADD CONSTRAINT signed_certifications_moves_id_fk
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for signed_certifications.personally_procured_move_id
ALTER TABLE signed_certifications
DROP CONSTRAINT signed_certifications_personally_procured_move_id_fkey;

ALTER TABLE signed_certifications
ADD CONSTRAINT signed_certifications_personally_procured_move_id_fkey
FOREIGN KEY (personally_procured_move_id)
REFERENCES personally_procured_moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for webhook_notifications.move_id
ALTER TABLE webhook_notifications
DROP CONSTRAINT webhook_notifications_move_id_fkey;

ALTER TABLE webhook_notifications
ADD CONSTRAINT webhook_notifications_move_id_fkey
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for mto_service_items.move_id
ALTER TABLE mto_service_items
DROP CONSTRAINT mto_service_items_move_id_fkey;

ALTER TABLE mto_service_items
ADD CONSTRAINT mto_service_items_move_id_fkey
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for payment_requests.move_id
ALTER TABLE payment_requests
DROP CONSTRAINT payment_requests_move_id_fkey;

ALTER TABLE payment_requests
ADD CONSTRAINT payment_requests_move_id_fkey
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for archived_move_documents.move_id
ALTER TABLE archived_move_documents
DROP CONSTRAINT archived_move_documents_move_id;

ALTER TABLE archived_move_documents
ADD CONSTRAINT archived_move_documents_move_id
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for archived_signed_certifications.move_id
ALTER TABLE archived_signed_certifications
DROP CONSTRAINT archived_signed_certifications_move_id;

ALTER TABLE archived_signed_certifications
ADD CONSTRAINT archived_signed_certifications_move_id
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for archived_personally_procured_moves.move_id
ALTER TABLE archived_personally_procured_moves
DROP CONSTRAINT archived_personally_procured_moves_move_id_fkey;

ALTER TABLE archived_personally_procured_moves
ADD CONSTRAINT archived_personally_procured_moves_move_id_fkey
FOREIGN KEY (move_id)
REFERENCES moves(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'mto_shipments' dependents
-------------------------------------------
-- Add ON DELETE CASCADE constraint for sit_extensions.mto_shipment_id
ALTER TABLE sit_extensions
DROP CONSTRAINT sit_extensions_mto_shipment_id_fkey;

ALTER TABLE sit_extensions
ADD CONSTRAINT sit_extensions_mto_shipment_id_fkey
FOREIGN KEY (mto_shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for mto_agents.mto_shipment_id
ALTER TABLE mto_agents
DROP CONSTRAINT mto_agents_mto_shipment_id_fkey;

ALTER TABLE mto_agents
ADD CONSTRAINT mto_agents_mto_shipment_id_fkey
FOREIGN KEY (mto_shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for reweighs.shipment_id
ALTER TABLE reweighs
DROP CONSTRAINT reweighs_shipment_id_fkey;

ALTER TABLE reweighs
ADD CONSTRAINT reweighs_shipment_id_fkey
FOREIGN KEY (shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for shipment_address_updates.shipment_id
ALTER TABLE shipment_address_updates
DROP CONSTRAINT shipment_address_updates_shipment_id_fkey;

ALTER TABLE shipment_address_updates
ADD CONSTRAINT shipment_address_updates_shipment_id_fkey
FOREIGN KEY (shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for evaluation_reports.shipment_id
ALTER TABLE evaluation_reports
DROP CONSTRAINT evaluation_reports_shipment_id_fkey;

ALTER TABLE evaluation_reports
ADD CONSTRAINT evaluation_reports_shipment_id_fkey
FOREIGN KEY (shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for mto_service_items.mto_shipment_id
ALTER TABLE mto_service_items
DROP CONSTRAINT mto_service_items_mto_shipment_id_fkey;

ALTER TABLE mto_service_items
ADD CONSTRAINT mto_service_items_mto_shipment_id_fkey
FOREIGN KEY (mto_shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for ppm_shipments.shipment_id
ALTER TABLE ppm_shipments
DROP CONSTRAINT ppm_shipment_mto_shipment_id_fkey;

ALTER TABLE ppm_shipments
ADD CONSTRAINT ppm_shipment_mto_shipment_id_fkey
FOREIGN KEY (shipment_id)
REFERENCES mto_shipments(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'ppm_shipments' dependents
-------------------------------------------
-- Add ON DELETE CASCADE constraint for signed_certifications.ppm_id
ALTER TABLE signed_certifications
ADD CONSTRAINT signed_certifications_ppm_shipments_id_fkey
FOREIGN KEY (ppm_id)
REFERENCES ppm_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for weight_tickets.ppm_shipment_id
ALTER TABLE weight_tickets
DROP CONSTRAINT weight_tickets_ppm_shipments_id_fkey;

ALTER TABLE weight_tickets
ADD CONSTRAINT weight_tickets_ppm_shipments_id_fkey
FOREIGN KEY (ppm_shipment_id)
REFERENCES ppm_shipments(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for progear_weight_tickets.ppm_shipment_id
ALTER TABLE progear_weight_tickets
DROP CONSTRAINT progear_weight_tickets_ppm_shipment_id_fkey;

ALTER TABLE progear_weight_tickets
ADD CONSTRAINT progear_weight_tickets_ppm_shipment_id_fkey
FOREIGN KEY (ppm_shipment_id)
REFERENCES ppm_shipments(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'mto_service_items' dependents
-------------------------------------------
-- Add ON DELETE CASCADE constraint for sit_address_updates.mto_service_item_id
ALTER TABLE sit_address_updates
DROP CONSTRAINT sit_address_updates_mto_service_item_id_fkey;

ALTER TABLE sit_address_updates
ADD CONSTRAINT sit_address_updates_mto_service_item_id_fkey
FOREIGN KEY (mto_service_item_id)
REFERENCES mto_service_items(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for service_items_customer_contacts.mtoservice_item_id
ALTER TABLE service_items_customer_contacts
DROP CONSTRAINT service_items_customer_contacts_mtoservice_item_id_fkey;

ALTER TABLE service_items_customer_contacts
ADD CONSTRAINT service_items_customer_contacts_mtoservice_item_id_fkey
FOREIGN KEY (mtoservice_item_id)
REFERENCES mto_service_items(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for payment_service_items.mto_service_item_id
ALTER TABLE payment_service_items
DROP CONSTRAINT payment_service_items_mto_service_item_id_fkey;

ALTER TABLE service_items_customer_contacts
ADD CONSTRAINT payment_service_items_mto_service_item_id_fkey
FOREIGN KEY (mtoservice_item_id)
REFERENCES mto_service_items(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for payment_service_item_params.payment_service_item_id
ALTER TABLE payment_service_item_params
DROP CONSTRAINT payment_service_item_params_payment_service_item_id_fkey;

ALTER TABLE payment_service_item_params
ADD CONSTRAINT payment_service_item_params_payment_service_item_id_fkey
FOREIGN KEY (payment_service_item_id)
REFERENCES payment_service_items(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for service_request_documents.mto_service_item_id
ALTER TABLE service_request_documents
DROP CONSTRAINT service_request_documents_mto_service_item_id_fkey;

ALTER TABLE service_request_documents
ADD CONSTRAINT service_request_documents_mto_service_item_id_fkey
FOREIGN KEY (mto_service_item_id)
REFERENCES mto_service_items(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for service_request_document_uploads.service_request_documents_id
ALTER TABLE service_request_document_uploads
DROP CONSTRAINT service_request_documents_service_request_documents_id_fkey;

ALTER TABLE service_request_document_uploads
ADD CONSTRAINT service_request_documents_service_request_documents_id_fkey
FOREIGN KEY (service_request_documents_id)
REFERENCES service_request_documents(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'payment_requests' dependents
-------------------------------------------
-- Add ON DELETE CASCADE constraint for payment_service_items.payment_request_id
ALTER TABLE payment_service_items
DROP CONSTRAINT payment_service_items_payment_request_id_fkey;

ALTER TABLE payment_service_items
ADD CONSTRAINT payment_service_items_payment_request_id_fkey
FOREIGN KEY (payment_request_id)
REFERENCES payment_requests(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for proof_of_service_docs.payment_request_id
ALTER TABLE proof_of_service_docs
DROP CONSTRAINT proof_of_service_docs_payment_request_id_fkey;

ALTER TABLE proof_of_service_docs
ADD CONSTRAINT proof_of_service_docs_payment_request_id_fkey
FOREIGN KEY (payment_request_id)
REFERENCES payment_requests(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for prime_uploads.proof_of_service_docs_id
ALTER TABLE proof_of_service_docs
DROP CONSTRAINT proof_of_service_docs_payment_request_id_fkey;

ALTER TABLE proof_of_service_docs
ADD CONSTRAINT proof_of_service_docs_payment_request_id_fkey
FOREIGN KEY (payment_request_id)
REFERENCES payment_requests(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for payment_request_to_interchange_control_numbers.payment_request_id
ALTER TABLE payment_request_to_interchange_control_numbers
DROP CONSTRAINT payment_request_to_icns_payment_request_id_fkey;

ALTER TABLE payment_request_to_interchange_control_numbers
ADD CONSTRAINT payment_request_to_icns_payment_request_id_fkey
FOREIGN KEY (payment_request_id)
REFERENCES payment_requests(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for edi_errors.payment_request_id
ALTER TABLE edi_errors
DROP CONSTRAINT edi_errors_payment_request_id_fkey;

ALTER TABLE edi_errors
ADD CONSTRAINT edi_errors_payment_request_id_fkey
FOREIGN KEY (payment_request_id)
REFERENCES payment_requests(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for edi_errors.interchange_control_number_id
ALTER TABLE edi_errors
DROP CONSTRAINT edi_errors_icn_id_fkey;

ALTER TABLE edi_errors
ADD CONSTRAINT edi_errors_icn_id_fkey
FOREIGN KEY (interchange_control_number_id)
REFERENCES payment_request_to_interchange_control_numbers(id)
ON DELETE CASCADE;

-------------------------------------------
-- 'archived_move_documents', 'archived_signed_certifications',
-- and 'archived_personally_procured_moves' dependents
-------------------------------------------

-- Add ON DELETE CASCADE constraint for archived_moving_expense_documents.move_document_id
ALTER TABLE archived_moving_expense_documents
DROP CONSTRAINT archived_moving_expense_documents_move_document_id_fkey;

ALTER TABLE archived_moving_expense_documents
ADD CONSTRAINT archived_moving_expense_documents_move_document_id_fkey
FOREIGN KEY (move_document_id)
REFERENCES archived_move_documents(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for archived_weight_ticket_set_documents.move_document_id
ALTER TABLE archived_weight_ticket_set_documents
DROP CONSTRAINT archived_weight_ticket_set_documents_move_document_id_fkey;

ALTER TABLE archived_weight_ticket_set_documents
ADD CONSTRAINT archived_weight_ticket_set_documents_move_document_id_fkey
FOREIGN KEY (move_document_id)
REFERENCES archived_move_documents(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for archived_move_documents.personally_procured_move_id
ALTER TABLE archived_move_documents
DROP CONSTRAINT archived_move_documents_personally_procured_move_id_fkey;

ALTER TABLE archived_move_documents
ADD CONSTRAINT archived_move_documents_personally_procured_move_id_fkey
FOREIGN KEY (personally_procured_move_id)
REFERENCES archived_personally_procured_moves(id)
ON DELETE CASCADE;

-- Add ON DELETE CASCADE constraint for archived_signed_certifications.personally_procured_move_id
ALTER TABLE archived_signed_certifications
DROP CONSTRAINT archived_signed_certifications_personally_procured_move_id_fkey;

ALTER TABLE archived_signed_certifications
ADD CONSTRAINT archived_signed_certifications_personally_procured_move_id_fkey
FOREIGN KEY (personally_procured_move_id)
REFERENCES archived_personally_procured_moves(id)
ON DELETE CASCADE;
