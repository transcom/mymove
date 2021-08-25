INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('9e70a393-aa8c-4ad2-a432-9b9bae71e7b8', (SELECT id FROM re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now())
