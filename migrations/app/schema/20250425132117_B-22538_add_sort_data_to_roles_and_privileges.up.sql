--B-22538   Michael Saki    Add updates for sort column in roles and privileges tables

UPDATE roles
SET sort = CASE role_type
    WHEN 'customer'                         THEN 1
    WHEN 'task_ordering_officer'            THEN 2
    WHEN 'task_invoicing_officer'           THEN 3
    WHEN 'contracting_officer'              THEN 4
    WHEN 'services_counselor'               THEN 5
    WHEN 'prime_simulator'                  THEN 6
    WHEN 'qae'                              THEN 7
    WHEN 'customer_service_representative'  THEN 8
    WHEN 'gsr'                              THEN 9
    WHEN 'headquarters'                     THEN 10
    ELSE sort
END
WHERE role_type IN (
    'customer',
    'task_ordering_officer',
    'task_invoicing_officer',
    'contracting_officer',
    'services_counselor',
    'prime_simulator',
    'qae',
    'customer_service_representative',
    'gsr',
    'headquarters'
);

UPDATE privileges
SET sort = CASE privilege_type
    WHEN 'supervisor'   THEN 1
    WHEN 'safety'       THEN 2
    ELSE sort
END
WHERE privilege_type IN ('supervisor', 'safety');