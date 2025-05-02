--B-22810   Pam Becker    Add new APPROVALS_REQUESTED shipment status
ALTER TYPE mto_shipment_status ADD VALUE IF NOT EXISTS 'APPROVALS_REQUESTED';