-- New enum for status
CREATE TYPE sit_address_status AS enum (
    'REQUESTED',
	'REJECTED',
    'APPROVED'
);

-- New table/columns
CREATE TABLE sit_address_updates
(
	id uuid PRIMARY KEY,
	mto_service_item_id uuid NOT NULL CONSTRAINT sit_address_updates_mto_service_item_id_fkey REFERENCES mto_service_items (id),
	old_address_id uuid NOT NULL CONSTRAINT sit_address_updates_old_address_id_fkey REFERENCES addresses (id),
	new_address_id uuid NOT NULL CONSTRAINT sit_address_updates_new_address_id_fkey REFERENCES addresses (id),
	status sit_address_status NOT NULL,
	distance int4 NOT NULL,
	reason text NOT NULL,
	contractor_remarks text NOT NULL,
	office_remarks text NULL
);

-- Column Comments
COMMENT on TABLE sit_address_updates IS 'Stores SIT destination address change requests for approval/rejection.';
COMMENT on COLUMN sit_address_updates.id IS 'uuid that represents this entity';
COMMENT on COLUMN sit_address_updates.mto_service_item_id IS 'Foreign key of the mto_service_items table';
COMMENT on COLUMN sit_address_updates.old_address_id IS 'Foreign key of addresses. Old address that will be replaced.';
COMMENT on COLUMN sit_address_updates.new_address_id IS 'Foreign key of addresses. New address that will replace the old address';
COMMENT on COLUMN sit_address_updates.status IS 'Current status of this request. Possible enum status(es): REQUESTED - Prime made this request and distance is greater than 50 miles, REJECTED - TXO rejected this request, APPROVED - TXO approved this request';
COMMENT on COLUMN sit_address_updates.distance IS 'The distance in miles between the old address and the new address. This is calculated and stored using the address zip codes.';
COMMENT on COLUMN sit_address_updates.reason IS 'A reason why this particular SIT address change is justified. TXOs would use the information here to accept or reject this SIT address change. Eg: Customer moving closer to family.';
COMMENT on COLUMN sit_address_updates.contractor_remarks IS 'Contractor remarks for the SIT address change. Eg: "Customer reached out to me this week & let me know they want to move closer to family."';
COMMENT on COLUMN sit_address_updates.office_remarks IS 'TXO remarks for the SIT address change.';
