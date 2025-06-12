-- B-22909 Daniel Jordan adding TERMINATED_FOR_CAUSE to mto_shipment_status enum
-- B-22810 Pam Becker    Add new APPROVALS_REQUESTED shipment status

ALTER TYPE mto_shipment_status ADD VALUE IF NOT EXISTS 'TERMINATED_FOR_CAUSE';
ALTER TYPE mto_shipment_status ADD VALUE IF NOT EXISTS 'APPROVALS_REQUESTED';