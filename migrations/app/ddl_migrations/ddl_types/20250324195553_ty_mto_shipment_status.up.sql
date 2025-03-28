-- B-22909 Daniel Jordan adding TERMINATED_FOR_CAUSE to mto_shipment_status enum
ALTER TYPE mto_shipment_status ADD VALUE IF NOT EXISTS 'TERMINATED_FOR_CAUSE';
