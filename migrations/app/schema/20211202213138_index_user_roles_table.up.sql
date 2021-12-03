-- Add a database migration to index the user_id and deleted_at columns of the users_roles table.
CREATE INDEX users_roles_user_id_idx ON users_roles(user_id);
CREATE INDEX users_roles_not_deleted_partial_idx ON users_roles (deleted_at)
    WHERE users_roles.deleted_at IS NULL;
COMMENT ON INDEX users_roles_not_deleted_partial_idx IS 'indexes user_roles that are not deleted';
