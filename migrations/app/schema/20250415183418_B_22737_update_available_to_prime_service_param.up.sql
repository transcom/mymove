UPDATE service_item_param_keys
    SET key = 'MTOEarliestRequestedPickup',
        description = 'Timestamp earliest non-PPM requested pickup date'
    WHERE key = 'MTOAvailableToPrimeAt';