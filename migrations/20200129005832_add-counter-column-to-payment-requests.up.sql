-- For now, make payment_request_number nullable.  Will need to figure out how to "fix"
-- existing payment_requests records.
alter table payment_requests
    add column payment_request_number text;

-- Can't do this yet
-- alter table payment_requests
--     add constraint payment_requests_number_unique_key unique (move_task_order_id, payment_request_number);
