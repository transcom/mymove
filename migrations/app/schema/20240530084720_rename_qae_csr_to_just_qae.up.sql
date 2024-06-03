-- Rename QAE/CSR to just QAE per E-05337
UPDATE roles
SET role_name = 'Quality Assurance Evaluator',
    role_type = 'qae',
    updated_at = now()
WHERE id = 'a2af3cc0-d0cd-4a29-8092-70ad45723090';