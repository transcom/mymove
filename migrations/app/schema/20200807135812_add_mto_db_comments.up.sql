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
COMMENT ON COLUMN orders.grade IS 'Customer''s rank. Should be found on their orders. Entered by the customer from a drop-down list. Includes ""civilian employee"';
COMMENT ON COLUMN orders.orders_type IS 'MilMove supports 4 orders types: Permanent change of station (PCS), Permanent change of assignment (PCA), retirement orders, and separation orders. In general, the moving process starts with the job/travel orders a customer receives from their service. In the orders, information describing rank, the duration of job/training, and their assigned location will determine if their entire dependent family can come, what the customer is allowed to bring, and how those items will arrive to their new location.';
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

COMMENT ON TABLE mto_shipments IS 'A move task order (MTO) shipment for a specific MTO.';

COMMENT ON TABLE mto_agents IS 'An agent is someone who can interact with movers on a customer''s behalf. There are receiving agents — people who can accept delivery at a location when the customer is not there. And releasing agents — people who can authorize a pickup from a location when the customer is not there.

Agents are assigned per shipment, not per move. The same person be an agent for multiple shipments. An agent is not a requirement for a shipment.';

COMMENT ON TABLE mto_service_items IS 'Service items associated with a particular MTO and shipment.';

COMMENT ON TABLE mto_service_item_dimensions IS 'The dimensions of a particular object within a particular MTO.';

COMMENT ON TABLE mto_service_item_customer_contacts IS 'Info related to a customer contact for a particular service item.';

COMMENT ON TABLE re_services IS 'Service codes allowed for a particular service type. Come from a finite list.';

COMMENT ON TABLE payment_requests IS 'The Prime can request payment for services rendered against an open MTO. Each request is tied to a particular service item.';

COMMENT ON TABLE payment_service_items IS 'A row in this table represents a service item from the MTO for which the prime wants to get paid.  It refers to the service item from the MTO but adds additional payment-specific data.';

COMMENT ON TABLE payment_service_item_params IS 'A row in this table represents an input parameter (i.e., a key/value pair) for computing the price of a service item. For example, the key might be "distance" and value "2000" (miles). A spreadsheet of inputs is available.';
