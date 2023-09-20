alter table users
add column okta_email text,
add column okta_id varchar DEFAULT '';

update users set okta_email = login_gov_email where login_gov_email is not NULL;

alter table users

drop column login_gov_email,
drop column login_gov_uuid;

CREATE INDEX ON users (okta_email);

COMMENT ON TABLE users IS 'Holds all users. Anyone who signs in to any of the mymove apps is automatically created in this table after signing in with okta.';
COMMENT ON COLUMN users.okta_id IS 'The okta id of the user.';
COMMENT ON COLUMN office_users.email IS 'The email of the office user. This will match their okta_email in the users table.';
