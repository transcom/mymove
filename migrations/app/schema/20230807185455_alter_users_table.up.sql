alter table users
add column  okta_email text,
add column okta_id varchar,
alter column login_gov_email DROP NOT NULL;