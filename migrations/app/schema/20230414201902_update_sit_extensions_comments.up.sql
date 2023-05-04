COMMENT ON TABLE "sit_extensions" IS 'Stores all the updates to SIT Durations that have been requested, and their details. Formerly known as SIT Extensions, SITDurationUpdates can include both increases and decreases to a SIT Duration.';
COMMENT ON COLUMN "sit_extensions"."mto_shipment_id" IS 'The MTO Shipment ID associated with this SIT Duration Update.';
COMMENT ON COLUMN "sit_extensions"."request_reason" IS 'One of a limited set of contractual reasons an Update to the SIT Duration can be requested.';
COMMENT ON COLUMN "sit_extensions"."contractor_remarks" IS 'Free form remarks from the contractor about this request to update the SIT Duration.';
COMMENT ON COLUMN "sit_extensions"."status" IS 'Status of this SIT Duration Update (Pending, Approved, or Denied).';
COMMENT ON COLUMN "sit_extensions"."approved_days" IS 'The number of days by which to update the SIT allowance. This number can be positive (increasing the SIT allowance) or negative (decreasing the SIT allowance)';
COMMENT ON COLUMN "sit_extensions"."decision_date" IS 'The date on which this request for a SIT Duration Update was approved or denied.';