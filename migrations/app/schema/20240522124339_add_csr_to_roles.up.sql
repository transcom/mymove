-- Insert CSR as a separate role as part of the feature to split QAE/CSR into two separate roles
INSERT INTO roles (id, role_type, created_at, updated_at, role_name)
VALUES (
        '72432922-BF2E-45DE-8837-1A458F5D1011',
        'customer_service_representative',
        now(),
        now(),
        'Customer Service Representative'
    );