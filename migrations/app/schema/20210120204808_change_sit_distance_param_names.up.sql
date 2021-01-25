-- Update a couple of key names/descriptions to better match their purpose.
UPDATE service_item_param_keys
SET key         = 'DistanceZipSITDest',
    description = 'Distance from shipment''s destination address to final destination address for delivery out of SIT',
    updated_at  = now()
WHERE key = 'DistanceZip5SITDest';

UPDATE service_item_param_keys
SET key         = 'DistanceZipSITOrigin',
    description = 'Distance from shipment''s original pickup address to actual pickup address for pickup out of SIT',
    updated_at  = now()
WHERE key = 'DistanceZip5SITOrigin'
