COMMENT ON COLUMN moves.status IS 'The current status of the Move. Allowed values are:
DRAFT, SUBMITTED, APPROVED, APPROVALS REQUESTED, CANCELED.
A Move''s status gets set to APPROVALS REQUESTED when the Prime creates new service
items. New service items can only be created if the Move''s previous status was
APPROVED or APPROVALS REQUESTED. The APPROVALS REQUESTED status lets the TOO know
they need to review the service items. As long as a Move has one or more service
items in SUBMITTED status (i.e. they have not been reviewed yet by the TOO), the
Move''s status remains APPROVALS REQUESTED. Once there aren''t any service items
in SUBMITTED status (i.e the TOO has either rejected or approved them), the
Move''s status goes back to APPROVED.';
