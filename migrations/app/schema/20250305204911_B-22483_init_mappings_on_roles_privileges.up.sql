-- B-22483   Ryan McHugh   update roles with static ids, add initial mappings to roles_privileges table

-- update roles with static ids
DO
$$
    DECLARE
        -- Current IDs for Role
        cust_id UUID;
        co_id UUID;
        tio_id UUID;
        too_id UUID;
        -- New IDs for Role
        cust_id_new UUID := 'c728caf3-5f9d-4db6-a9d1-7cd8ff013b2e';
        co_id_new UUID := '5496a188-69dc-4ae4-9dab-ce6c063d648f';
        tio_id_new UUID := 'c19a5d5f-d320-4972-b294-1d760ee4b899';
        too_id_new UUID := '2b21e867-78c3-4980-95a1-c8242b78baba';
    BEGIN
        select id into cust_id from roles where role_type = 'customer';
        select id into co_id from roles where role_type = 'contracting_officer';
        select id into tio_id from roles where role_type = 'task_invoicing_officer';
        select id into too_id from roles where role_type = 'task_ordering_officer';

        IF cust_id <> cust_id_new THEN
            update roles set role_type = 'customer1' where id = cust_id;
            insert into roles values (cust_id_new,'customer',now(),now(),'Customer');
            update users_roles set role_id = cust_id_new where role_id = cust_id;
            delete from roles where id = cust_id;
        END IF;

        IF co_id <> co_id_new THEN
            update roles set role_type = 'contracting_officer1' where id = co_id;
            insert into roles values (co_id_new,'contracting_officer',now(),now(),'Contracting Officer');
            update users_roles set role_id = co_id_new where role_id = co_id;
            delete from roles where id = co_id;
        END IF;

        IF tio_id <> tio_id_new THEN
            update roles set role_type = 'task_invoicing_officer1' where id = tio_id;
            insert into roles values (tio_id_new,'task_invoicing_officer',now(),now(),'Task Invoicing Officer');
            update users_roles set role_id = tio_id_new where role_id = tio_id;
            delete from roles where id = tio_id;
        END IF;

        IF too_id <> too_id_new THEN
            update roles set role_type = 'task_ordering_officer1' where id = too_id;
            insert into roles values (too_id_new,'task_ordering_officer',now(),now(),'Task Ordering Officer');
            update users_roles set role_id = too_id_new where role_id = too_id;
            delete from roles where id = too_id;
        END IF;
END;
$$;

-- insert applicable roles for supervisor privilege
INSERT INTO roles_privileges (id, role_id, privilege_id, created_at, updated_at) VALUES
    ('c261ca4c-185b-4a69-adfb-a8cce625d93e'::uuid,'c728caf3-5f9d-4db6-a9d1-7cd8ff013b2e'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Customer
    ('e20c1e39-fea8-43ef-a56c-7e276bb33c87'::uuid,'2b21e867-78c3-4980-95a1-c8242b78baba'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Task Ordering Officer
    ('2d32e971-1350-46e9-931c-1a53b2e3af3b'::uuid,'c19a5d5f-d320-4972-b294-1d760ee4b899'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Task Invoicing Officer
    ('66eaee36-311c-4184-aeb9-178b1d407c0a'::uuid,'5496a188-69dc-4ae4-9dab-ce6c063d648f'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Contracting Officer
    ('afe69e74-0abd-4ccf-9a15-e727fe18c746'::uuid,'010bdae1-8ebe-44c9-b8ee-8c4477fae2a6'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Services Counselor
    ('14665c1a-888f-4969-9dd8-21077b41f401'::uuid,'63c07db0-5a7d-499c-ab64-90c08f74f654'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Prime Simulator
    ('f33326bf-e56d-4bb8-999b-5961662be70b'::uuid,'a2af3cc0-d0cd-4a29-8092-70ad45723090'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Quality Assurance Evaluator
    ('e62984fc-03d8-496e-8ba1-e4c979120f4a'::uuid,'72432922-bf2e-45de-8837-1a458f5d1011'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Customer Service Representative
    ('e0d49f7a-96a8-4bde-874c-83bb322558f1'::uuid,'20d7deea-4010-424e-9f64-714a46e18c3c'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Government Surveillance Representative
    ('1ed7973d-c022-48fa-9534-12278cb9d98c'::uuid,'0da36914-fcc1-4965-b49c-b4a0d447514c'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()); -- Headquarters

-- insert applicable roles for safety privilege
INSERT INTO roles_privileges (id, role_id, privilege_id, created_at, updated_at) VALUES
    ('dac8f0e3-6e0e-476e-848b-d12b663319ea'::uuid,'2b21e867-78c3-4980-95a1-c8242b78baba'::uuid,'43f77473-2ecd-4b06-920a-e1e003f63c18'::uuid, now(), now()), -- Task Ordering Officer
    ('091aca53-4732-46d7-977f-fa4b7b4f67b4'::uuid,'c19a5d5f-d320-4972-b294-1d760ee4b899'::uuid,'43f77473-2ecd-4b06-920a-e1e003f63c18'::uuid, now(), now()), -- Task Invoicing Officer
    ('0b4d838a-bd5c-4921-8140-344c8917c80b'::uuid,'010bdae1-8ebe-44c9-b8ee-8c4477fae2a6'::uuid,'43f77473-2ecd-4b06-920a-e1e003f63c18'::uuid, now(), now()), -- Services Counselor
    ('911b83a8-b5b2-41fe-8d94-91e212818837'::uuid,'a2af3cc0-d0cd-4a29-8092-70ad45723090'::uuid,'43f77473-2ecd-4b06-920a-e1e003f63c18'::uuid, now(), now()), -- Quality Assurance Evaluator
    ('ce985179-bd80-429e-96b7-9a8f1ca41ab6'::uuid,'72432922-bf2e-45de-8837-1a458f5d1011'::uuid,'43f77473-2ecd-4b06-920a-e1e003f63c18'::uuid, now(), now()), -- Customer Service Representative
    ('91915f91-a031-495b-9b25-5699bf24be5b'::uuid,'0da36914-fcc1-4965-b49c-b4a0d447514c'::uuid,'43f77473-2ecd-4b06-920a-e1e003f63c18'::uuid, now(), now()); -- Headquarters