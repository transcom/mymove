-- Instead of using ID here, "Key" is unique so we can safely use that
-- B-22663 Switching NTSPackingFactor from PRICING origin to SYSTEM origin
-- This is because pricing origin params are generated AFTER pricing has completed,
-- but we need this param available DURING pricing as it is needed for pricing math
update service_item_param_keys set origin = 'SYSTEM' where key = 'NTSPackingFactor';