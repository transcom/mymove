ALTER TABLE re_services ADD COLUMN priority integer NOT NULL DEFAULT 99;

UPDATE re_services SET priority = 1 WHERE id = '50f1179a-3b72-4fa1-a951-fe5bcc70bd14' AND name = 'Dom. Destination Price';
UPDATE re_services SET priority = 1 WHERE id = '8d600f25-1def-422d-b159-617c7d59156e' AND name = 'Dom. Linehaul';
UPDATE re_services SET priority = 1 WHERE id = '2bc3e5cb-adef-46b1-bde9-55570bfdd43e' AND name = 'Dom. Origin Price';
UPDATE re_services SET priority = 1 WHERE id = 'bdea5a8d-f15f-47d2-85c9-bba5694802ce' AND name = 'Dom. Packing';
UPDATE re_services SET priority = 1 WHERE id = '4b85962e-25d3-4485-b43c-2497c4365598' AND name = 'Dom. Shorthaul';
UPDATE re_services SET priority = 1 WHERE id = '07051352-4715-49b5-88e7-045b7541919d' AND name = 'Int''l. C->O Shipping & LH';
UPDATE re_services SET priority = 1 WHERE id = '67ba1eaf-6ffd-49de-9a69-497be7789877' AND name = 'Int''l. HHG Pack';
UPDATE re_services SET priority = 1 WHERE id = 'd0bb2cae-838a-4fc7-8efc-f7c6ad57431d' AND name = 'Int''l. O->C Shipping & LH';
UPDATE re_services SET priority = 1 WHERE id = '56bb94cd-f160-4239-a028-b31ffc641eb7' AND name = 'Int''l. O->O Shipping & LH';