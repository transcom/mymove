ALTER TABLE re_service_items ADD sort int;

COMMENT ON COLUMN re_service_items.sort IS 'Sort order for service items to be displayed for a given shipment type.';

update re_service_items set sort = 1 where service_id in (select id from re_services where code in ('ISLH','UBP'));
update re_service_items set sort = 2 where service_id in (select id from re_services where code = 'POEFSC');
update re_service_items set sort = 3 where service_id in (select id from re_services where code = 'PODFSC');
update re_service_items set sort = 4 where service_id in (select id from re_services where code in ('IHPK','IUBPK'));
update re_service_items set sort = 5 where service_id in (select id from re_services where code in ('IHUPK','IUBUPK'));