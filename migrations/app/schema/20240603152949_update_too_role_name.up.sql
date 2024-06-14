UPDATE
    roles
SET
    role_name = 'Task Ordering Officer', role_type = 'task_ordering_officer'
WHERE
    role_name = 'Transportation Ordering Officer';

COMMENT ON COLUMN mto_shipments.rejection_reason IS 'Not currently used, until the "reject" or "cancel" a shipment feature is implemented. When the Task Ordering Officer rejects or cancels a shipment, they will explain why';
COMMENT ON COLUMN mto_shipments.approved_date IS 'The date when the Task Ordering Officer approves the shipment, and it is added to the Move Task Order for the Prime contractor';
COMMENT ON COLUMN orders.tac IS '(For HHG shipments) Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Task Ordering Officer adds this information to the MTO.';
COMMENT ON COLUMN orders.nts_tac IS '(For NTS shipments) Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Task Ordering Officer adds this information to the MTO.';
COMMENT ON COLUMN roles.role_type IS 'The name of the role in snake case. Current values are: ''task_ordering_officer'', ''transportation_invoicing_officer'', ''customer'', ''ppm_office_users'', ''contracting_officer''.';
