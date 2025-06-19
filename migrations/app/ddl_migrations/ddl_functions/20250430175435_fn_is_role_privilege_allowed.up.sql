--B-22539 Elizabeth Perkins Add function to check role privilege association

CREATE OR REPLACE FUNCTION is_role_privilege_allowed(input_role TEXT, input_privilege TEXT)
 RETURNS boolean
 LANGUAGE plpgsql
AS $$
BEGIN
    IF EXISTS(
        SELECT 1
        FROM roles_privileges
        JOIN roles ON roles_privileges.role_id = roles.id
        JOIN privileges ON roles_privileges.privilege_id = privileges.id
        WHERE roles.role_type = input_role AND privileges.privilege_type = input_privilege
    ) THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$;