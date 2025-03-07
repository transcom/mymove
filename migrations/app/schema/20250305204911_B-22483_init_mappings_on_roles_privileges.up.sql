-- B-22483 - add initial mappings to roles_privileges table

-- insert applicable roles for supervisor privilege
INSERT INTO roles_privileges (id, role_id, privilege_id, created_at, updated_at) VALUES
    ('e20c1e39-fea8-43ef-a56c-7e276bb33c87'::uuid,'2b21e867-78c3-4980-95a1-c8242b78baba'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Task Ordering Officer
    ('2d32e971-1350-46e9-931c-1a53b2e3af3b'::uuid,'c19a5d5f-d320-4972-b294-1d760ee4b899'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Task Invoicing Officer
    ('66eaee36-311c-4184-aeb9-178b1d407c0a'::uuid,'5496a188-69dc-4ae4-9dab-ce6c063d648f'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Contracting Officer
    ('afe69e74-0abd-4ccf-9a15-e727fe18c746'::uuid,'010bdae1-8ebe-44c9-b8ee-8c4477fae2a6'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Services Counselor
    ('14665c1a-888f-4969-9dd8-21077b41f401'::uuid,'63c07db0-5a7d-499c-ab64-90c08f74f654'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Prime Simulator
    ('f33326bf-e56d-4bb8-999b-5961662be70b'::uuid,'a2af3cc0-d0cd-4a29-8092-70ad45723090'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Quality Assurance Evaluator
    ('e62984fc-03d8-496e-8ba1-e4c979120f4a'::uuid,'72432922-bf2e-45de-8837-1a458f5d1011'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Customer Service Representative
    ('e0d49f7a-96a8-4bde-874c-83bb322558f1'::uuid,'20d7deea-4010-424e-9f64-714a46e18c3c'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()), -- Government Surveillance Representative
    ('1ed7973d-c022-48fa-9534-12278cb9d98c'::uuid,'0da36914-fcc1-4965-b49c-b4a0d447514c'::uuid,'463c2034-d197-4d9a-897e-8bbe64893a31'::uuid, now(), now()); -- Headquarters



