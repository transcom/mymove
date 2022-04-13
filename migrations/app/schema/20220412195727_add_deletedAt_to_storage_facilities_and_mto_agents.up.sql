ALTER TABLE mto_agents
ADD deleted_at timestamp with time zone;

ALTER TABLE storage_facilities
ADD deleted_at timestamp with time zone;

COMMENT ON COLUMN mto_agents.deleted_at IS 'Indicates whether the mto agent has been soft deleted or not, and when it was soft deleted.';
COMMENT ON COLUMN storage_facilities.deleted_at IS 'Indicates whether the storage facility has been soft deleted or not, and when it was soft deleted.';
