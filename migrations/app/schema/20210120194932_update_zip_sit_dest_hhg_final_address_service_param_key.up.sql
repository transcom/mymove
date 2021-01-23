-- Rename ZipSITAddress to ZipSITDestHHGFinalAddress and update description
UPDATE service_item_param_keys
  SET key = 'ZipSITDestHHGFinalAddress',
    description = 'SIT Final Destination Address ZIP from MTOServiceItem',
    updated_at = now()
  WHERE key = 'ZipSITAddress';

-- Update ZipDestAddress description to clarify its source is the MTOShipment
UPDATE service_item_param_keys
  SET description = 'Destination address ZIP from MTOShipment', updated_at = now()
  WHERE key = 'ZipDestAddress';
