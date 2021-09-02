CREATE TYPE sit_extension_request_reason AS enum (
    'SERIOUS_ILLNESS_MEMBER',
    'SERIOUS_ILLNESS_DEPENDENT',
    'IMPENDING_ASSIGNEMENT',
    'DIRECTED_TEMPORARY_DUTY',
    'NONAVAILABILITY_OF_CIVILIAN_HOUSING',
    'AWAITING_COMPLETION_OF_RESIDENCE',
    'OTHER'
);

COMMENT ON TYPE sit_extension_request_reason IS 'List of reasons a SIT extension can be requested for';

CREATE TYPE sit_extension_status AS enum (
    'PENDING',
    'APPROVED',
    'DENIED'
);

COMMENT ON TYPE sit_extension_status IS 'List of possible statuses for a SIT Extension';

CREATE TABLE sit_extensions (
  id uuid NOT NULL primary key,
  mto_shipment_id uuid NOT NULL
    CONSTRAINT sit_extensions_mto_shipment_id_fkey
      REFERENCES mto_shipments,
  request_reason sit_extension_request_reason NOT NULL,
  contractor_remarks varchar,
  requested_days int NOT NULL,
  status sit_extension_status NOT NULL,
  approved_days int,
  decision_date timestamp,
  office_remarks varchar,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);
CREATE INDEX on sit_extensions (mto_shipment_id);

COMMENT on TABLE sit_extensions IS 'Stores all the sit extensions that have been requested, and their details.';
COMMENT on COLUMN sit_extensions.mto_shipment_id IS 'MTO Shipment ID associated with this SIT Extension.';
COMMENT on COLUMN sit_extensions.request_reason IS 'One of a limited set of contractual reasons an extension can be requested.';
COMMENT on COLUMN sit_extensions.contractor_remarks IS 'Free form remarks from the contractor about the extension request.';
COMMENT on COLUMN sit_extensions.requested_days IS 'Number of requested days to extend the SIT allowance by.';
COMMENT on COLUMN sit_extensions.status IS 'Status of this SIT Extension request (Pending, Approved, or Denied).';
COMMENT on COLUMN sit_extensions.approved_days IS 'Number of days approved to extend the SIT allowance by. May differ from the original requested days.';
COMMENT on COLUMN sit_extensions.decision_date IS 'Date that the TOO approved or deined this extension request.';
COMMENT on COLUMN sit_extensions.office_remarks IS 'Any comments from the TOO on the approval or denial of this request.';
