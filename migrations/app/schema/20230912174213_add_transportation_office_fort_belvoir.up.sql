-- Insert Fort Belvoir Transportation Office
-- Jira ticket: MB-17020

insert into addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
values
    ('96e22159-d829-4f68-b53e-44f790dc06df', '10109 Gridley Rd', 'Bldg 314' 'Fort Belvoir', 'VA', '22060', now(), now());

insert into transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
values
    ('a877a317-be5f-482b-a126-c91a34be9290' , 'PPPO Fort Belvoir - USA', 'BGAC', '96e22159-d829-4f68-b53e-44f790dc06df', 38.6837701, -77.1399064, now(), now());

insert into office_phone_lines
    (id, transportation_office_id, number, created_at, updated_at)
values
    ('ee3c08c7-ca09-4dba-9e4b-b570e7916432', 'a877a317-be5f-482b-a126-c91a34be9290', '800-521-9959', now(), now());

insert into office_emails
    (id, transportation_office_id, email, created_at, updated_at)
values
    ('4a6d1f99-b38f-4378-850f-d72c43c67299', 'a877a317-be5f-482b-a126-c91a34be9290', 'usarmy.belvoir.asc.mbx.jppsoma@army.mil', now(), now());
