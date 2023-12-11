ALTER TABLE sit_extensions
ADD COLUMN members_expense BOOLEAN DEFAULT FALSE;
COMMENT on TABLE sit_extensions.members_expense IS 'Whether or not the service member is responsible for expenses of SIT (i.e. if SIT extension request was denied).';