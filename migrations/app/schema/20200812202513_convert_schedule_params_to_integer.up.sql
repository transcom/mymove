UPDATE service_item_param_keys
SET type = 'INTEGER'
WHERE key = 'ServicesScheduleOrigin'
OR key = 'ServicesScheduleDest'
OR key = 'SITScheduleOrigin'
OR key = 'SITScheduleDest';