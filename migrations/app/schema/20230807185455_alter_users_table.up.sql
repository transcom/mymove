-- commenting these out due to introduction of these columns at an earlier migration file during change from login.gov to Okta, this is no longer needed

alter table users
add column okta_email text,
add column okta_id varchar;
