-- Change FSCPriceDifferenceInCents to be a DECIMAL instead of an INTEGER since fuel prices can have tenths of cents
UPDATE service_item_param_keys
SET type = 'DECIMAL'
WHERE key = 'FSCPriceDifferenceInCents';
