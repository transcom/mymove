COMMENT ON TABLE move_task_orders IS 'Contains all the information on the Move Task Order (MTO). There is one MTO per a customer''s move.';
COMMENT ON COLUMN move_task_orders.is_canceled IS 'Identifies if an MTO was canceled.';
COMMENT ON COLUMN move_task_orders.created_at IS 'Date & time the MTO was created';
COMMENT ON COLUMN move_task_orders.updated_at IS 'Date & time the MTO was last updated';
COMMENT ON COLUMN move_task_orders.ppm_type IS 'Identifies whether a move is a full PPM or a partial PPM — the customer moving everything or only some things.';
COMMENT ON COLUMN move_task_orders.ppm_estimated_weight IS 'Estimated weight of the part of a customer''s belongings that they will move in a PPM. Unit is pounds. Customer does the estimation for PPMs.';
COMMENT ON COLUMN move_task_orders.contractor_id IS 'Unique identifier for the prime contractor.';
COMMENT ON COLUMN move_task_orders.available_to_prime_at IS 'Date & time the TOO made the MTO available to the prime contractor.';

COMMENT ON TABLE orders IS 'A customer''s move is initiated or changed based on orders issued to them by their service. Details change based on the service, but for MilMove purposes the orders will include what type of orders they are, where the customer is going, when the customer needs to get there, and other info.';
COMMENT ON COLUMN orders.service_member_id IS 'Unique identifier for the customer — the person who has the orders and is moving.';
COMMENT ON COLUMN orders.origin_duty_station_id IS 'Unique identifier for the duty station the customer is moving from. Not the same as the text version of the name.';
COMMENT ON COLUMN orders.new_duty_station_id IS 'Unique identifier for the duty station the customer being assigned to. Not the same as the text version of the name.';
COMMENT ON COLUMN orders.orders_number IS 'Information found on the customer''s orders assigned by their service that uniquely identifies the document. Entered in MilMove by the counselor or TOO.';
COMMENT ON COLUMN orders.grade IS 'Customer''s rank. Should be found on their orders. Entered by the customer from a drop-down list. Includes "civilian employee"';
COMMENT ON COLUMN orders.orders_type IS 'MilMove supports 4 orders types: Permanent change of station (PCS), Permanent change of assignment (PCA), retirement orders, and separation orders.
In general, the moving process starts with the job/travel orders a customer receives from their service. In the orders, information describing rank, the duration of job/training, and their assigned location will determine if their entire dependent family can come, what the customer is allowed to bring, and how those items will arrive to their new location.';
COMMENT ON COLUMN orders.orders_type_detail IS 'Selected from a drop-down list. Includes more specific info about the kind of move orders the customer received.
List includes:
- Shipment of HHG permitted - PCS with TDY en route
- Shipment of HHG restricted or prohibited
- HHG restricted area-HHG prohibited
- Course of instruction 20 weeks or more
- Shipment of HHG prohibited but authorized within 20 weeks
- Delayed approval 20 weeks or more';
COMMENT ON COLUMN orders.issue_date IS 'Date on which the customer''s orders were issued by their branch of service.';
COMMENT ON COLUMN orders.report_by_date IS 'Date by which the customer must report to their new duty station or assignment.';
COMMENT ON COLUMN orders.tac IS 'Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Transportation Ordering Officer adds this information to the MTO.';

COMMENT ON TABLE service_members IS 'The customer is the service member or other individual who received move orders. They are the point person for the move, though family members are likely coming with them.';
COMMENT ON COLUMN service_members.edipi IS 'The customer''s Department of Defense ID number, which is used as their unique ID.';
COMMENT ON COLUMN service_members.first_name IS 'The customer''s first name';
COMMENT ON COLUMN service_members.last_name IS 'The customer''s last name';
COMMENT ON COLUMN service_members.affiliation IS 'The customer''s branch of service';
COMMENT ON COLUMN service_members.personal_email IS 'The customer''s email address';
COMMENT ON COLUMN service_members.telephone IS 'The customer''s phone number';

COMMENT ON TABLE entitlements IS 'Service members are entitled to have the government pay to move a certain amount of weight, based on their rank, whether or not they have dependents, and whether their destination is CONUS or OCONUS. "Entitlements" is an older term, and the services now call these "allowances".';
COMMENT ON COLUMN entitlements.dependents_authorized IS 'A yes/no field reflecting whether dependents are authorized on the customer''s move orders';
COMMENT ON COLUMN entitlements.total_dependents IS 'An integer reflecting the total number of dependents that are authorized on the customer''s move orders. For UB shipments, if dependents are authorized, each dependent adds to the customer''s weight allowance. (Note that the exact amount depends on the dependent''s age.)';
COMMENT ON COLUMN entitlements.non_temporary_storage IS 'A yes/no field reflecting whether the customer is requesting a Non Temporary Storage shipment';
COMMENT ON COLUMN entitlements.privately_owned_vehicle IS 'A yes/no field reflecting whether the customer has a privately owned vehicle that will need to be shipped as part of their move';
COMMENT ON COLUMN entitlements.storage_in_transit IS 'The maximum number of days of storage in transit allowed by the customer''s move orders';
COMMENT ON COLUMN entitlements.authorized_weight IS 'The maximum number of pounds the Prime contractor is authorized to move for the customer';

COMMENT ON TABLE mto_shipments IS 'A move task order (MTO) shipment for a specific MTO.';
COMMENT ON COLUMN mto_shipments.scheduled_pickup_date IS 'The scheduled pickup date the Prime contractor schedules for a shipment in consultation with the customer';
COMMENT ON COLUMN mto_shipments.requested_pickup_date IS 'The date the customer is requesting that a given shipment be picked up.';
COMMENT ON COLUMN mto_shipments.customer_remarks IS 'The remarks field is where the customer can describe special circumstances for their shipment, in order to inform the Prime contractor of any unique shipping and handling needs.';
COMMENT ON COLUMN mto_shipments.pickup_address_id IS 'The customer''s pickup address for a shipment';
COMMENT ON COLUMN mto_shipments.destination_address_id IS 'The customer''s destination address for a shipment';
COMMENT ON COLUMN mto_shipments.prime_estimated_weight IS 'The estimated weight of a shipment, provided by the Prime contractor after they survey a customer''s shipment';
COMMENT ON COLUMN mto_shipments.prime_estimated_weight_recorded_date IS 'Date when the Prime contractor records the shipment''s estimated weight';
COMMENT ON COLUMN mto_shipments.prime_actual_weight IS 'The actual weight of a shipment, provided by the Prime contractor after they pack, pickup, and weigh a customer''s shipment';
COMMENT ON COLUMN mto_shipments.shipment_type IS 'The type of shipment. The list includes:
1. Personally procured move (PPM)
2. Household goods move (HHG)
3. Non-temporary storage (NTS)';
COMMENT ON COLUMN mto_shipments.status IS 'The status of a shipment. The list of statuses includes:
1. New
2. Move approved
3. Approvals requested
4. Payment requested
5. Move complete';
COMMENT ON COLUMN mto_shipments.rejection_reason IS 'Not currently used, until the "reject" or "cancel" a shipment feature is implemented. When the Transportation Ordering Officer rejects or cancels a shipment, they will explain why';
COMMENT ON COLUMN mto_shipments.actual_pickup_date IS 'The actual pickup date when the Prime contractor picks up the customer''s shipment';
COMMENT ON COLUMN mto_shipments.approved_date IS 'The date when the Transportation Ordering Officer approves the shipment, and it is added to the Move Task Order for the Prime contractor';

COMMENT ON TABLE mto_agents IS 'An agent is someone who can interact with movers on a customer''s behalf. There are receiving agents — people who can accept delivery at a location when the customer is not there. And releasing agents — people who can authorize a pickup from a location when the customer is not there.
Agents are assigned per shipment, not per move. The same person be an agent for multiple shipments. An agent is not a requirement for a shipment.';
COMMENT ON COLUMN mto_agents.mto_shipment_id IS 'The shipment this particular agent applies to — a receiving agent for one shipment is not necessarily an agent for other shipments.';
COMMENT ON COLUMN mto_agents.agent_type IS 'Either RELEASING agent, or RECEIVING agent. Someone who can authorize a pickup, or who can authorize a delivery.';
COMMENT ON COLUMN mto_agents.first_name IS 'First name of the agent, not the customer.';
COMMENT ON COLUMN mto_agents.last_name IS 'Last name of the agent, not the customer.';
COMMENT ON COLUMN mto_agents.email IS 'Email contact for the agent.';
COMMENT ON COLUMN mto_agents.phone IS 'Phone number for the agent.';

COMMENT ON TABLE mto_service_items IS 'Service items associated with a particular MTO and shipment.';
COMMENT ON COLUMN mto_service_items.reason IS 'A reason why this particular service item is justified. TXOs would use the information here to accept or reject a service item (crating, shuttling, etc.).';
COMMENT ON COLUMN mto_service_items.pickup_postal_code IS 'ZIP for the location where the shipment pickup is taking place.';
COMMENT ON COLUMN mto_service_items.description IS 'Description of the item that the service item applies to. If it''s a request for crating, for example, this describes the item being crated (piano, moose head, etc.).';

COMMENT ON TABLE mto_service_item_dimensions IS 'The dimensions of a particular object within a particular MTO.';
COMMENT ON COLUMN mto_service_item_dimensions.type IS 'Identifies if the dimensions apply to the item being crated, or to the crate itself. (ITEM or CRATE)';
COMMENT ON COLUMN mto_service_item_dimensions.length_thousandth_inches IS 'Length in thousandth inches. 1000 thou = 1 inch.';
COMMENT ON COLUMN mto_service_item_dimensions.height_thousandth_inches IS 'Height in thousandth inches. 1000 thou = 1 inch.';
COMMENT ON COLUMN mto_service_item_dimensions.width_thousandth_inches IS 'Width in thousandth inches. 1000 thou = 1 inch.';

COMMENT ON TABLE mto_service_item_customer_contacts IS 'Info related to a customer contact for a particular service item.';
COMMENT ON COLUMN mto_service_item_customer_contacts.type IS 'Is this the FIRST or SECOND method of contact for a customer when communicating about this MTO service item.';
COMMENT ON COLUMN mto_service_item_customer_contacts.time_military IS 'Time that the service item was delivered.';
COMMENT ON COLUMN mto_service_item_customer_contacts.first_available_delivery_date IS 'First available date that Prime can deliver this service item.';

COMMENT ON TABLE re_services IS 'Service codes allowed for a particular service type. Come from a finite list.';
COMMENT ON COLUMN re_services.code IS 'Specific code for a service type. Examples: DDDSIT, DPK';
COMMENT ON COLUMN re_services.name IS 'Full human-readable name of the service type. Examples: Domestic destination SIT delivery, Domestic packing';

COMMENT ON TABLE payment_requests IS 'The Prime can request payment for services rendered against an open MTO. Each request is tied to a particular service item.';
COMMENT ON COLUMN payment_requests.rejection_reason IS 'Entered by TIO to explain why they rejected a payment request.';
COMMENT ON COLUMN payment_requests.move_task_order_id IS 'Which MTO the payment request should be applied to.';
COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID';

COMMENT ON TABLE payment_service_items IS 'A row in this table represents a service item from the MTO for which the prime wants to get paid.  It refers to the service item from the MTO but adds additional payment-specific data.';
COMMENT ON COLUMN payment_service_items.status IS 'Can be REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID';
COMMENT ON COLUMN payment_service_items.price_cents IS 'Price of the service item, measured in cents (US).';
COMMENT ON COLUMN payment_service_items.rejection_reason IS 'Why the TIO rejected payment for this service item.';

COMMENT ON TABLE payment_service_item_params IS 'A row in this table represents an input parameter (i.e., a key/value pair) for computing the price of a service item. For example, the key might be "distance" and value "2000" (miles). A spreadsheet of inputs is available.';
