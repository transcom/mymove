ALTER TABLE sit_extensions
ADD COLUMN members_expense BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN sit_extensions.members_expense IS 'Denotes that the TOO rejected this extension request AND converted it to member''s expense (could be used in MTO view/history to show exactly when a shipment was converted)';

ALTER TABLE mto_shipments
ADD COLUMN members_expense BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN mto_shipments.members_expense IS 'Whether or not the service member is responsible for expenses of SIT (i.e. if SIT extension request was denied).';