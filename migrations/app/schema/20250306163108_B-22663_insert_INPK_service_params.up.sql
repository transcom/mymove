-- inserting INPK NTS market factor as an optional param on IHPK
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
    ('A9995EA1-5F53-43CA-834E-C5D2220179CE'::uuid, '67ba1eaf-6ffd-49de-9a69-497be7789877', 'b2b433ae-18b5-4362-9f74-50771175a838', NOW(), NOW(), true); -- nts packing factor
