CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Finds office users who have never signed in so we must create a user record
-- for them that can be updated on first sign in
WITH office_never_signed_in AS (
    SELECT *
    FROM office_users
             LEFT JOIN users ON office_users.email = users.login_gov_email
    WHERE users.login_gov_email IS NULL
)
INSERT INTO users (id, login_gov_email, created_at, updated_at)
SELECT uuid_generate_v4(), email, now(), now()
FROM office_never_signed_in;

-- Now that we've created the user record, update the office_user record with
-- the new user.id
WITH office_associate_user AS (
    SELECT users.*
    FROM office_users
    JOIN users ON office_users.email = users.login_gov_email
    WHERE office_users.user_id IS NULL
)
UPDATE office_users
SET user_id = office_associate_user.id
FROM office_associate_user
WHERE office_users.email = office_associate_user.login_gov_email;

-- Finds admin users who have never signed in so we must create a user record
-- for them that can be updated on first sign in
WITH admin_never_signed_in AS (
    SELECT *
    FROM admin_users
             LEFT JOIN users ON admin_users.email = users.login_gov_email
    WHERE users.login_gov_email IS NULL
)
INSERT INTO users (id, login_gov_email, created_at, updated_at)
SELECT uuid_generate_v4(), email, now(), now()
FROM admin_never_signed_in;

-- Now that we've created the user record, update the admin_user record with
-- the new user.id
WITH admin_associate_user AS (
    SELECT users.*
    FROM admin_users
    JOIN users ON admin_users.email = users.login_gov_email
    WHERE admin_users.user_id IS NULL
)
UPDATE admin_users
SET user_id = admin_associate_user.id
FROM admin_associate_user
WHERE admin_users.email = admin_associate_user.login_gov_email;