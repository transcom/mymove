ALTER TYPE mto_shipment_status ADD VALUE 'CANCELLATION_REQUESTED';
COMMENT ON COLUMN mto_shipments.status IS 'The status of a shipment. The list of statuses includes:
1. DRAFT
2. SUBMITTED
3. APPROVED
4. REJECTED
5. CANCELLATION_REQUESTED';
