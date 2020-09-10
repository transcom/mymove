-- client_certs
COMMENT ON TABLE client_certs IS 'Holds the SSL/TLS certificates authorized for MilMove and indicates to which parts of the app they have access.';
COMMENT ON COLUMN client_certs.sha256_digest IS 'The encrypted signature of the certificate';
COMMENT ON COLUMN client_certs.subject IS 'The entity the certificate belongs to';
COMMENT ON COLUMN client_certs.allow_dps_auth_api IS 'Indicates whether or not the cert grants access to the DPS Auth API';
COMMENT ON COLUMN client_certs.allow_orders_api IS 'Indicates whether or not the cert grants access to the Orders API';
COMMENT ON COLUMN client_certs.allow_air_force_orders_read IS 'Indicates whether or not the cert grants view-only access to Air Force orders';
COMMENT ON COLUMN client_certs.allow_air_force_orders_write IS 'Indicates whether or not the cert grants edit access to Air Force orders';
COMMENT ON COLUMN client_certs.allow_army_orders_read IS 'Indicates whether or not the cert grants view-only access to Army orders';
COMMENT ON COLUMN client_certs.allow_army_orders_write IS 'Indicates whether or not the cert grants edit access to Army orders';
COMMENT ON COLUMN client_certs.allow_coast_guard_orders_read IS 'Indicates whether or not the cert grants view-only access to Coast Guard orders';
COMMENT ON COLUMN client_certs.allow_coast_guard_orders_write IS 'Indicates whether or not the cert grants edit access to Coast Guard orders';
COMMENT ON COLUMN client_certs.allow_marine_corps_orders_read IS 'Indicates whether or not the cert grants view-only access to Marine Corps orders';
COMMENT ON COLUMN client_certs.allow_marine_corps_orders_write IS 'Indicates whether or not the cert grants edit access to Marine Corps orders';
COMMENT ON COLUMN client_certs.allow_navy_orders_read IS 'Indicates whether or not the cert grants view-only access to Navy orders';
COMMENT ON COLUMN client_certs.allow_navy_orders_write IS 'Indicates whether or not the cert grants edit access to Navy orders';
COMMENT ON COLUMN client_certs.allow_prime IS 'Indicates whether or not the cert grants access to the Prime API';
COMMENT ON COLUMN client_certs.created_at IS 'Date & time the client cert was created';
COMMENT ON COLUMN client_certs.updated_at IS 'Date & time the client cert was last updated';

-- electronic_orders
COMMENT ON TABLE electronic_orders IS 'Represents the electronic move orders issued by a particular branch of the military';
COMMENT ON COLUMN electronic_orders.orders_number IS 'A (generally) unique number identifying the orders, corresponding to the ORDERS number (Army), the CT SDN (Navy, Marines), the SPECIAL ORDER NO (Air Force), the Travel Order No. (Coast Guard), or the Travel Authorization Number (Civilian)';
COMMENT ON COLUMN electronic_orders.edipi IS 'Electronic Data Interchange Personal Identifier, the 10 digit DoD ID Number of the service member';
COMMENT ON COLUMN electronic_orders.issuer IS 'The organization that issued the orders (army, navy, etc.)';
COMMENT ON COLUMN electronic_orders.created_at IS 'Date & time the electronic order was created';
COMMENT ON COLUMN electronic_orders.updated_at IS 'Date & time the electronic order was last updated';

-- electronic_orders_revisions
COMMENT ON TABLE electronic_orders_revisions IS 'Represents revisions or edits to an issued set of electronic move orders';
COMMENT ON COLUMN electronic_orders_revisions.electronic_order_id IS 'The UUID of the electronic orders being revised';
COMMENT ON COLUMN electronic_orders_revisions.seq_num IS 'An integer representing the sequence number for this revision. As orders are amended, the revision with the highest sequence number is considered the current, authoritative version of the orders, even if date_issued is earlier.';
COMMENT ON COLUMN electronic_orders_revisions.given_name IS 'The first name of the service member';
COMMENT ON COLUMN electronic_orders_revisions.middle_name IS 'The middle name or initial of the service member';
COMMENT ON COLUMN electronic_orders_revisions.family_name IS 'The last name of the service member';
COMMENT ON COLUMN electronic_orders_revisions.name_suffix IS 'The suffix of the service member''s name, if any (Jr. Sr., III, etc.)';
COMMENT ON COLUMN electronic_orders_revisions.affiliation IS 'The service member''s affiliated military branch (army, navy, etc.)';
COMMENT ON COLUMN electronic_orders_revisions.paygrade IS 'The DoD paygrade or rank of the service member';
COMMENT ON COLUMN electronic_orders_revisions.title IS 'The preferred form of address for the service member. Used mainly for ranks that have multiple possible titles.';
COMMENT ON COLUMN electronic_orders_revisions.status IS 'Indicates whether these Orders are authorized, RFO (Request For Orders), or canceled';
COMMENT ON COLUMN electronic_orders_revisions.date_issued IS 'The date and time thath these orders were cut. If omitted, the current date and time will be used.';
COMMENT ON COLUMN electronic_orders_revisions.no_cost_move IS 'If true, indicates that these orders do not authorize any move expenses. If false, these orders are a PCS and should authorize expenses.';
COMMENT ON COLUMN electronic_orders_revisions.tdy_en_route IS 'TDY (Temporary Duty Yonder) en-route. If omitted, assume false.';
COMMENT ON COLUMN electronic_orders_revisions.tour_type IS 'Accompanied or Unaccompanied - indicates whether or not dependents are authorized to accompany the service member on the move. If omitted, assume accompanied.';
COMMENT ON COLUMN electronic_orders_revisions.orders_type IS 'The type of orders for this move (joining the military, retirement, training, etc.)';
COMMENT ON COLUMN electronic_orders_revisions.has_dependents IS 'Indicates whether or not the service member has any dependents (spouse, children, etc.)';
COMMENT ON COLUMN electronic_orders_revisions.losing_uic IS 'The Unit Identification Code for the unit the service member is moving away from. A six character code that identifies each DoD entity.';
COMMENT ON COLUMN electronic_orders_revisions.losing_unit_name IS 'The human-readable name of the losing unit';
COMMENT ON COLUMN electronic_orders_revisions.losing_unit_city IS 'The city of the losing unit. May be FPO or APO for OCONUS commands.';
COMMENT ON COLUMN electronic_orders_revisions.losing_unit_locality IS 'The locality of the losing unit. Will be the state for US units.';
COMMENT ON COLUMN electronic_orders_revisions.losing_unit_country IS 'The ISO 3166-1 alpha-2 country code for the losing unit. If blank, but city, locality, or postal_code are not blank, assume US';
COMMENT ON COLUMN electronic_orders_revisions.losing_unit_postal_code IS 'The postal code of the losing unit. Will be the ZIP code for US units.';
COMMENT ON COLUMN electronic_orders_revisions.gaining_uic IS 'The Unit Identification Code for the unit the service member is moving to. May be blank if these are separation orders.';
COMMENT ON COLUMN electronic_orders_revisions.gaining_unit_name IS 'The human-readable name of the gaining unit';
COMMENT ON COLUMN electronic_orders_revisions.gaining_unit_city IS 'The city of the gaining unit. May be FPO or APO for OCONUS commands.';
COMMENT ON COLUMN electronic_orders_revisions.gaining_unit_locality IS 'The locality of the gaining unit. Will be the state for US units.';
COMMENT ON COLUMN electronic_orders_revisions.gaining_unit_country IS 'The ISO 3166-1 alpha-2 country code for the gaining unit. If blank, but city, locality, or postal_code are not blank, assume US';
COMMENT ON COLUMN electronic_orders_revisions.gaining_unit_postal_code IS 'The postal code of the gaining unit. Will be the ZIP code for US units.';
COMMENT ON COLUMN electronic_orders_revisions.report_no_earlier_than IS 'Earliest date that the service member is allowed to report for duty at the new duty station. If omitted, the member is allowed to report as early as desired.';
COMMENT ON COLUMN electronic_orders_revisions.report_no_later_than IS 'Latest date that the service member is allowed to report for duty at the new duty station. Should be included for most Orders types, but can be missing for Separation / Retirement Orders.';
COMMENT ON COLUMN electronic_orders_revisions.hhg_tac IS 'Transportation Account Code. Used for accounting purposes in HHG expenses.';
COMMENT ON COLUMN electronic_orders_revisions.hhg_sdn IS 'Standard Document Number. Used for routing money for an HHG expense.';
COMMENT ON COLUMN electronic_orders_revisions.hhg_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for HHG expenses.';
COMMENT ON COLUMN electronic_orders_revisions.nts_tac IS 'Transportation Account Code. Used for accounting purposes in NTS expenses.';
COMMENT ON COLUMN electronic_orders_revisions.nts_sdn IS 'Standard Document Number. Used for routing money for an NTS expense.';
COMMENT ON COLUMN electronic_orders_revisions.nts_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for NTS expenses.';
COMMENT ON COLUMN electronic_orders_revisions.pov_shipment_tac IS 'Transportation Account Code. Used for accounting purposes in POV shipment expenses.';
COMMENT ON COLUMN electronic_orders_revisions.pov_shipment_sdn IS 'Standard Document Number. Used for routing money for a POV shipment expense.';
COMMENT ON COLUMN electronic_orders_revisions.pov_shipment_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for POV shipment expenses.';
COMMENT ON COLUMN electronic_orders_revisions.pov_storage_tac IS 'Transportation Account Code. Used for accounting purposes in POV storage expenses.';
COMMENT ON COLUMN electronic_orders_revisions.pov_storage_sdn IS 'Standard Document Number. Used for routing money for a POV storage expense.';
COMMENT ON COLUMN electronic_orders_revisions.pov_storage_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for POV storage expenses.';
COMMENT ON COLUMN electronic_orders_revisions.ub_tac IS 'Transportation Account Code. Used for accounting purposes in UB expenses.';
COMMENT ON COLUMN electronic_orders_revisions.ub_sdn IS 'Standard Document Number. Used for routing money for a UB expense.';
COMMENT ON COLUMN electronic_orders_revisions.ub_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for UB expenses.';
COMMENT ON COLUMN electronic_orders_revisions.comments IS 'Free-form text that may or may not contain information relevant to moving';
COMMENT ON COLUMN electronic_orders_revisions.created_at IS 'Date & time the revision for the electronic order was created';
COMMENT ON COLUMN electronic_orders_revisions.updated_at IS 'Date & time the revision for the electronic order was last updated';

-- entitlements
COMMENT ON COLUMN entitlements.created_at IS 'Date & time the entitlement was created';
COMMENT ON COLUMN entitlements.updated_at IS 'Date & time the entitlement was last updated';

-- move_documents
COMMENT ON TABLE move_documents IS 'Represents files that have been uploaded to document steps in the move process';
COMMENT ON COLUMN move_documents.move_id IS 'The UUID of the move this document was uploaded for';
COMMENT ON COLUMN move_documents.document_id IS 'The UUID of the base document record';
COMMENT ON COLUMN move_documents.move_document_type IS 'The type of document this is (storage expense, shipment summary, etc.)';
COMMENT ON COLUMN move_documents.status IS 'The status of this document in the review process. Can be:
1. AWAITING_REVIEW
2. HAS_ISSUE
3. EXCLUDE_FROM_CALCULATION
4. OK';
COMMENT ON COLUMN move_documents.notes IS 'Any pertinent info on the document that needs to be noted for the review process';
COMMENT ON COLUMN move_documents.title IS 'The title of the document';
COMMENT ON COLUMN move_documents.personally_procured_move_id IS 'The UUID of the PPM move this document was uploaded for';
COMMENT ON COLUMN move_documents.created_at IS 'Date & time the move document was created';
COMMENT ON COLUMN move_documents.updated_at IS 'Date & time the move document was last updated';
COMMENT ON COLUMN move_documents.deleted_at IS 'Date & time the move document was deleted';

-- moves
COMMENT ON COLUMN moves.orders_id IS 'Unique identifier for the orders issued for this move.';

-- moving_expense_documents
COMMENT ON TABLE moving_expense_documents IS 'Used for PPM - represents the files the Service Member submits to document their move expenses';
COMMENT ON COLUMN moving_expense_documents.move_document_id IS 'The UUID for the move document this record refers to';
COMMENT ON COLUMN moving_expense_documents.moving_expense_type IS 'The type of expense (gas, tolls, rental equipment, etc.)';
COMMENT ON COLUMN moving_expense_documents.requested_amount_cents IS 'The amount for which the service member is requesting reimbursement, in cents';
COMMENT ON COLUMN moving_expense_documents.payment_method IS 'The payment method for this expense';
COMMENT ON COLUMN moving_expense_documents.receipt_missing IS 'Indicates whether or not the receipt is missing from this request';
COMMENT ON COLUMN moving_expense_documents.storage_start_date IS 'For storage expenses, indicates the date when the item first entered storage';
COMMENT ON COLUMN moving_expense_documents.storage_end_date IS 'For storage expenses, indicates the date when the item left storage';
COMMENT ON COLUMN moving_expense_documents.created_at IS 'Date & time the moving expense document was created';
COMMENT ON COLUMN moving_expense_documents.updated_at IS 'Date & time the moving expense document was last updated';
COMMENT ON COLUMN moving_expense_documents.deleted_at IS 'Date & time the moving expense document was deleted';

-- mto_agents - Correcting a typo on the table comment
COMMENT ON TABLE mto_agents IS 'An agent is someone who can interact with movers on a customer''s behalf. There are receiving agents — people who can accept delivery at a location when the customer is not there. And releasing agents — people who can authorize a pickup from a location when the customer is not there.
Agents are assigned per shipment, not per move. The same person may be an agent for multiple shipments. An agent is not a requirement for a shipment.';
COMMENT ON COLUMN mto_agents.created_at IS 'Date & time the agent was created';
COMMENT ON COLUMN mto_agents.updated_at IS 'Date & time the agent was last updated';

-- mto_service_item_customer_contacts
COMMENT ON TABLE mto_service_item_customer_contacts IS 'Holds the data for when the Prime contacted the customer to deliver their shipment but were unable to do so. Used to justify the Prime putting the shipment into a SIT facility.';
COMMENT ON COLUMN mto_service_item_customer_contacts.mto_service_item_id IS 'The UUID of the SIT service item this customer contact justifies';
COMMENT ON COLUMN mto_service_item_customer_contacts.type IS 'Either the FIRST or SECOND attempt at contacting the customer for delivery';
COMMENT ON COLUMN mto_service_item_customer_contacts.time_military IS 'The time the Prime contacted the customer, in military format (HHMMZ)';
COMMENT ON COLUMN mto_service_item_customer_contacts.first_available_delivery_date IS 'The date when the Prime attempted to deliver the shipment';
COMMENT ON COLUMN mto_service_item_customer_contacts.created_at IS 'Date & time the customer contact was created';
COMMENT ON COLUMN mto_service_item_customer_contacts.updated_at IS 'Date & time the customer contact was last updated';

-- mto_service_item_dimensions
COMMENT ON COLUMN mto_service_item_dimensions.mto_service_item_id IS 'The UUID of the service item these dimensions are associated with';
COMMENT ON COLUMN mto_service_item_dimensions.created_at IS 'Date & time the service item dimension was created';
COMMENT ON COLUMN mto_service_item_dimensions.updated_at IS 'Date & time the service item dimension was last updated';

-- mto_service_items
COMMENT ON COLUMN mto_service_items.move_id IS 'The UUID of the move this service item is for';
COMMENT ON COLUMN mto_service_items.mto_shipment_id IS 'The UUID of the shipment this service item is for - optional';
COMMENT ON COLUMN mto_service_items.re_service_id IS 'The UUID of the service code for this service item';
COMMENT ON COLUMN mto_service_items.rejection_reason IS 'The reason why the TOO might have rejected this service item request';
COMMENT ON COLUMN mto_service_items.status IS 'The status of this service item in the review process. Can be:
1. SUBMITTED
2. APPROVED
3. REJECTED';
COMMENT ON COLUMN mto_service_items.approved_at IS 'Date & time the service item was marked as APPROVED';
COMMENT ON COLUMN mto_service_items.rejected_at IS 'Date & time the service item was marked as REJECTED';
COMMENT ON COLUMN mto_service_items.created_at IS 'Date & time the service item was created';
COMMENT ON COLUMN mto_service_items.updated_at IS 'Date & time the service item was last updated';

-- mto_shipments
COMMENT ON COLUMN mto_shipments.move_id IS 'The UUID of the move this shipment is for';
COMMENT ON COLUMN mto_shipments.scheduled_pickup_date IS 'The pickup date the Prime contractor schedules for a shipment after consultation with the customer';
COMMENT ON COLUMN mto_shipments.secondary_pickup_address_id IS 'The secondary pickup address for this shipment';
COMMENT ON COLUMN mto_shipments.secondary_delivery_address_id IS 'The secondary delivery address for this shipment';
COMMENT ON COLUMN mto_shipments.status IS 'The status of a shipment. The list of statuses includes:
1. DRAFT
2. SUBMITTED
3. APPROVED
4. REJECTED';
COMMENT ON COLUMN mto_shipments.first_available_delivery_date IS 'Date the Prime provides to the customer so the customer can plan their own travel accordingly. We need to collect the FADD on the MTO so there is a record of what the Prime said they told the customer in case a situation arises in which the customer is unavailable to receive delivery of a shipment and the Prime wants to put the shipment in SIT.';
COMMENT ON COLUMN mto_shipments.required_delivery_date IS 'Latest date the Prime can deliver a customer''s shipment without violating the contract. RDD is the last date in the spread of available dates calculated from the scheduled pickup date.';
COMMENT ON COLUMN mto_shipments.days_in_storage IS 'Related specifically to SIT. Total number of days a shipment was in temporary storage, determined after it comes out of SIT.';
COMMENT ON COLUMN mto_shipments.requested_delivery_date IS 'Entered by the customer. Available delivery dates and required delivery date are calculated based on this date. Not at all a guarantee that this is the date the Prime will deliver the shipment.';
COMMENT ON COLUMN mto_shipments.distance IS 'Distance the shipment traveled, in miles';
COMMENT ON COLUMN mto_shipments.created_at IS 'Date & time the shipment was created';
COMMENT ON COLUMN mto_shipments.updated_at IS 'Date & time the shipment was last updated';

-- webhook_notifications
COMMENT ON TABLE webhook_notifications IS 'Represents the notifications that will be sent to an external client about changes in our database. Used to notify the Prime when Prime-available moves and related objects have been updated.';
COMMENT ON COLUMN webhook_notifications.event_key IS 'A string used to identify which object this notification pertains to (PaymentRequest, MTOShipment, etc.), and how it was modified (PaymentRequest.Create, MTOShipment.Update, etc.)';
COMMENT ON COLUMN webhook_notifications.trace_id IS 'The UUID representing the specific transaction this notification represents';
COMMENT ON COLUMN webhook_notifications.move_id IS 'The UUID for the move affected by this change';
COMMENT ON COLUMN webhook_notifications.object_id IS 'The UUID for the specific object that was modified';
COMMENT ON COLUMN webhook_notifications.payload IS 'A JSON payload containing the updates for this object';
COMMENT ON COLUMN webhook_notifications.status IS 'The status of this notification. Can be:
1. PENDING
2. SENT
3. FAILED';
COMMENT ON COLUMN webhook_notifications.created_at IS 'Date & time the webhook notification was created';
COMMENT ON COLUMN webhook_notifications.updated_at IS 'Date & time the webhook notification was last updated';

-- webhook_subscriptions
COMMENT ON TABLE webhook_subscriptions IS 'Represents subscribers who expect certain notifications to be pushed to their servers. Used for the Prime and Prime-related events specifically.';
COMMENT ON COLUMN webhook_subscriptions.subscriber_id IS 'The UUID of the Prime contractor this subscription belongs to';
COMMENT ON COLUMN webhook_subscriptions.status IS 'The status of this subscription. Can be:
1. ACTIVE
2. DISABLED';
COMMENT ON COLUMN webhook_subscriptions.event_key IS 'A string used to represent which events this subscriber expects to be notified about. Corresponds to the possible `event_key` values in `webhook_notifications`';
COMMENT ON COLUMN webhook_subscriptions.callback_url IS 'The URL to which the notifications for this subscription will be pushed to';
COMMENT ON COLUMN webhook_subscriptions.created_at IS 'Date & time the webhook subscription was created';
COMMENT ON COLUMN webhook_subscriptions.updated_at IS 'Date & time the webhook subscription was last updated';