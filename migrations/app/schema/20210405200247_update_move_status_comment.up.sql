COMMENT ON COLUMN moves.status IS 'The current status of the Move. Allowed values are:
DRAFT,
SUBMITTED,
APPROVED,
APPROVALS REQUESTED,
CANCELED,
NEEDS SERVICE COUNSELING,
SERVICE COUNSELING COMPLETED.

A Move starts out in DRAFT status upon creation. When the customer submits the
move, it either gets routed to Service Counseling (based on the origin duty
station), at which point its status becomes NEEDS SERVICE COUNSELING, or gets
submitted directly to the TOO and its status is SUBMITTED.

Once a service counselor completes their work, the status becomes
SERVICE COUNSELING COMPLETED, and this move now appears in the TOO''s queue.
When a TOO approves this move, its status becomes APPROVED, just like moves that
were sent straight to the TOO as SUBMITTED.

A Move''s status gets set to APPROVALS REQUESTED when the Prime creates new service
items. New service items can only be created if the Move''s previous status was
APPROVED or APPROVALS REQUESTED. The APPROVALS REQUESTED status lets the TOO know
they need to review the service items. As long as a Move has one or more service
items in SUBMITTED status (i.e. they have not been reviewed yet by the TOO), the
Move''s status remains APPROVALS REQUESTED. Once there aren''t any service items
in SUBMITTED status (i.e the TOO has either rejected or approved them), the
Move''s status goes back to APPROVED.';
