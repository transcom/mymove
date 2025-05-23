alter table backup_contacts
add column if not exists first_name text,
add column if not exists last_name text;

comment on column backup_contacts.first_name is 'First name of the backup contact';
comment on column backup_contacts.last_name is 'Last name of the backup contact';

select distinct ltrim(name) name, substr(name,1,position(' ' in name)) first_name, substr(name,position(' ' in name)+1,255) last_name
from backup_contacts
order by ltrim(name);

update backup_contacts
   set first_name = substr(name,1,position(' ' in name)-1),
   	   last_name = substr(name,position(' ' in name)+1,255);

--alter table backup_contacts
--drop column if exists name;
