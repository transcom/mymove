UPDATE service_item_param_keys
    SET type = 'DECIMAL'
    WHERE key IN ('CubicFeetCrating', 'CubicFeetBilled');
