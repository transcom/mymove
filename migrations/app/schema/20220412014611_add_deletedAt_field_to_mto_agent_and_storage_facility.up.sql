ALTER TABLE mto_agents
ADD deleted_at timestamp with time zone;

ALTER TABLE storage_facilities
ADD deleted_at timestamp with time zone;
