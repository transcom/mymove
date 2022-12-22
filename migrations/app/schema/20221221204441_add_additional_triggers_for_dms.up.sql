SELECT add_audit_history_table(
    target_table := 'customer_support_remarks',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'distance_calculations',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'documents',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'duty_location_names',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'duty_locations',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'edi_errors',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'edi_processings',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'edi_processings',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'electronic_orders',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'electronic_orders_revisions',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'entitlements',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'evaluation_reports',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'fuel_eia_diesel_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'ghc_diesel_fuel_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'ghc_domestic_transit_times',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'invoice_number_trackers',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'invoices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'jppso_region_state_assignments',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'jppso_regions',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'moving_expenses',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'notifications',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'office_emails',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'office_emails',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'office_phone_lines',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'office_users',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'organizations',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'payment_request_to_interchange_control_numbers',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'payment_service_item_params',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'payment_service_items',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'personally_procured_moves',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'postal_code_to_gblocs',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'ppm_shipments',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'prime_uploads',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'progear_weight_tickets',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'pws_violations',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_contract_years',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_contracts',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_domestic_accessorial_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_domestic_linehaul_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_domestic_other_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_domestic_service_area_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_domestic_service_areas',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_intl_accessorial_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_intl_accessorial_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_intl_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_rate_areas',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_services',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_shipment_type_prices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_task_order_fees',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_zip3s',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 're_zip5_rate_areas',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'report_violations',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'roles',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'service_item_param_keys',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'service_params',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'sit_extensions',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'transportation_accounting_codes',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'transportation_offices',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'transportation_service_provider_performances',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'transportation_service_providers',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'uploads',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'users',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'users_roles',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'weight_tickets',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);

SELECT add_audit_history_table(
    target_table := 'zip3_distances',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at'
    ]
);
