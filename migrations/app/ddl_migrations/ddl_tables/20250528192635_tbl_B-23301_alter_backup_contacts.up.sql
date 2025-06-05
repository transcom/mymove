-- B-23301  Jim Hawks  added columns for first name and last name

alter table backup_contacts
  add column if not exists first_name text default '',
  add column if not exists last_name text default '';

comment on column backup_contacts.first_name is 'First name of the backup contact';
comment on column backup_contacts.last_name is 'Last name of the backup contact';

alter table backup_contacts
alter column name drop not null;


