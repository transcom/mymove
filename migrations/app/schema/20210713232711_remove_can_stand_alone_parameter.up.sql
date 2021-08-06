-- delete payment_service_item_params tied to CanStandAlone
delete from payment_service_item_params where service_item_param_key_id = (select id from service_item_param_keys where key = 'CanStandAlone');
-- delete service_item_param_keys tied to CanStandAlone
delete from service_params where service_item_param_key_id = (select id from service_item_param_keys where key = 'CanStandAlone');
-- remove the service_item_param_keys
delete from service_item_param_keys where key = 'CanStandAlone';
