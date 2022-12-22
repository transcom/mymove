-- Delete all user role records for office users with the ppm_office_user role
DELETE FROM users_roles WHERE role_id = (SELECT id FROM roles WHERE role_type = 'ppm_office_users');

DELETE FROM roles WHERE role_type = 'ppm_office_users';