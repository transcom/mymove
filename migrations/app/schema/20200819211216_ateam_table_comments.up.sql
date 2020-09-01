COMMENT ON TABLE fuel_eia_diesel_prices IS 'Stores SDDC Fuel Surcharge rate information; used by pre-GHC HHG moves.';
COMMENT ON COLUMN fuel_eia_diesel_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN fuel_eia_diesel_prices.pub_date IS 'The date this rate was published.';
COMMENT ON COLUMN fuel_eia_diesel_prices.rate_start_date IS 'The start date that this rate is applicable (inclusive).';
COMMENT ON COLUMN fuel_eia_diesel_prices.rate_end_date IS 'The end date that this rate is applicable (inclusive).';
COMMENT ON COLUMN fuel_eia_diesel_prices.eia_price_per_gallon_millicents IS 'The national average price per gallon in millicents for this period as determined by the EIA (Energy Information Administration).';
COMMENT ON COLUMN fuel_eia_diesel_prices.baseline_rate IS 'The calculated baseline fuel surcharge rate in cents for this period.';
COMMENT ON COLUMN fuel_eia_diesel_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN fuel_eia_diesel_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE ghc_diesel_fuel_prices IS 'Represents the weekly average diesel fuel price; used in GHC pricing.';
COMMENT ON COLUMN ghc_diesel_fuel_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN ghc_diesel_fuel_prices.fuel_price_in_millicents IS 'The national average price per gallon in millicents for the week following the publication date as determined by the EIA (Energy Information Administration).';
COMMENT ON COLUMN ghc_diesel_fuel_prices.publication_date IS 'The date this rate was published.';
COMMENT ON COLUMN ghc_diesel_fuel_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN ghc_diesel_fuel_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE invoice_number_trackers IS 'Tracks latest sequence numbers in SCAC/year groupings; this sequence number is part of the generated invoice number.';
COMMENT ON COLUMN invoice_number_trackers.standard_carrier_alpha_code IS 'The associated SCAC for this sequence number (see the transportation_service_providers table).';
COMMENT ON COLUMN invoice_number_trackers.year IS 'The associated year for this sequence number.';
COMMENT ON COLUMN invoice_number_trackers.sequence_number IS 'The last used sequence number for the given SCAC/year.';

COMMENT ON TABLE invoices IS 'Represents an invoice sent to GEX; only used by pre-GHC HHG moves at the moment.';
COMMENT ON COLUMN invoices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN invoices.status IS 'Status of this invoice; options are DRAFT, IN_PROCESS, SUBMITTED, SUBMISSION_FAILURE, UPDATE_FAILURE.';
COMMENT ON COLUMN invoices.invoiced_date IS 'Timestamp when this invoice was sent to GEX.';
COMMENT ON COLUMN invoices.invoice_number IS 'A unique invoice number. Format is SCAC + two digit year + sequence number (with a suffix of -01, -02, etc. appended for subsequent invoices on the same shipment).';
COMMENT ON COLUMN invoices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN invoices.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN invoices.approver_id IS 'The office user that approved this invoice.';
COMMENT ON COLUMN invoices.user_uploads_id IS 'The associated uploads used as justification for this invoice.';

COMMENT ON TABLE payment_requests IS 'Represents a payment request from the GHC prime contractor.';
COMMENT ON COLUMN payment_requests.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN payment_requests.is_final IS 'True if this is the final payment request for the move task order (MTO).';
COMMENT ON COLUMN payment_requests.rejection_reason IS 'The reason the payment request was rejected (if it was rejected).';
COMMENT ON COLUMN payment_requests.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN payment_requests.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN payment_requests.move_id IS 'The associated move for the payment request.';
COMMENT ON COLUMN payment_requests.status IS 'The status of the payment request; options are PENDING, REVIEWED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID.';
COMMENT ON COLUMN payment_requests.requested_at IS 'Timestamp when the payment request was requested.';
COMMENT ON COLUMN payment_requests.reviewed_at IS 'Timestamp when the payment request was reviewed.';
COMMENT ON COLUMN payment_requests.sent_to_gex_at IS 'Timestamp when the payment request was sent to GEX.';
COMMENT ON COLUMN payment_requests.received_by_gex_at IS 'Timestamp when the payment request was received by GEX.';
COMMENT ON COLUMN payment_requests.paid_at IS 'Timestamp when the payment request was paid.';
COMMENT ON COLUMN payment_requests.payment_request_number IS 'A human-readable identifier for the payment request; format is <reference_id>-<sequence_number>.';
COMMENT ON COLUMN payment_requests.sequence_number IS 'The sequence number of this payment request for the associated move (the first payment request would be 1).';

COMMENT ON TABLE payment_service_item_params IS 'Represents the parameters (key/value pairs) for a given service item in a payment request.';
COMMENT ON COLUMN payment_service_item_params.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN payment_service_item_params.payment_service_item_id IS 'The associated service item in the payment request.';
COMMENT ON COLUMN payment_service_item_params.service_item_param_key_id IS 'The key for this parameter.';
COMMENT ON COLUMN payment_service_item_params.value IS 'The value for this parameter.';
COMMENT ON COLUMN payment_service_item_params.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN payment_service_item_params.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE payment_service_items IS 'Represents the service items associated with a given payment request.';
COMMENT ON COLUMN payment_service_items.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN payment_service_items.payment_request_id IS 'The associated payment request.';
COMMENT ON COLUMN payment_service_items.status IS 'The payment status of this service item; options are REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID.';
COMMENT ON COLUMN payment_service_items.price_cents IS 'The calculated price in cents for this service item (as determined by the GHC rate engine).';
COMMENT ON COLUMN payment_service_items.rejection_reason IS 'The reason payment for a service item was rejected (if it was rejected).';
COMMENT ON COLUMN payment_service_items.requested_at IS 'Timestamp when payment for the service item was requested.';
COMMENT ON COLUMN payment_service_items.approved_at IS 'Timestamp when payment for the service item was approved.';
COMMENT ON COLUMN payment_service_items.denied_at IS 'Timestamp when payment for the service item was denied.';
COMMENT ON COLUMN payment_service_items.sent_to_gex_at IS 'Timestamp when payment for the service item was sent to GEX.';
COMMENT ON COLUMN payment_service_items.paid_at IS 'Timestamp when payment for the service item was paid.';
COMMENT ON COLUMN payment_service_items.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN payment_service_items.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN payment_service_items.mto_service_item_id IS 'The associated MTO service item for which payment is requested.';

COMMENT ON TABLE prime_uploads IS 'Represents uploads made by the GHC prime contractor.';
COMMENT ON COLUMN prime_uploads.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN prime_uploads.proof_of_service_docs_id IS 'The associated set of proof of service documents this upload belongs to.';
COMMENT ON COLUMN prime_uploads.contractor_id IS 'The associated contractor for this upload.';
COMMENT ON COLUMN prime_uploads.upload_id IS 'The associated set of metadata for this upload.';
COMMENT ON COLUMN prime_uploads.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN prime_uploads.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN prime_uploads.deleted_at IS 'Timestamp when the upload was deleted.';

COMMENT ON TABLE proof_of_service_docs IS 'Ties together a set of uploads as proof of service documents.';
COMMENT ON COLUMN proof_of_service_docs.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN proof_of_service_docs.payment_request_id IS 'The associated payment request that these proof of service documents support.';
COMMENT ON COLUMN proof_of_service_docs.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN proof_of_service_docs.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_contract_years IS 'Represents the "years" included in a GHC pricing contract (see sheet 5b).';
COMMENT ON COLUMN re_contract_years.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_contract_years.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_contract_years.name IS 'The name of this contract year (e.g., "Base Period Year 1").';
COMMENT ON COLUMN re_contract_years.start_date IS 'The start date for this contract year (inclusive).';
COMMENT ON COLUMN re_contract_years.end_date IS 'The end date for this contract year (inclusive).';
COMMENT ON COLUMN re_contract_years.escalation IS 'The escalation factor for this specific contract year.';
COMMENT ON COLUMN re_contract_years.escalation_compounded IS 'The compounded escalation factor after applying previous year''s escalations.';
COMMENT ON COLUMN re_contract_years.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_contract_years.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_contracts IS 'Represents a GHC pricing contract; helps to tie together all data in that contract.';
COMMENT ON COLUMN re_contracts.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_contracts.code IS 'A short, human-readable code that uniquely identifies a contract.';
COMMENT ON COLUMN re_contracts.name IS 'A longer, more descriptive name for the contract.';
COMMENT ON COLUMN re_contracts.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_contracts.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_domestic_accessorial_prices IS 'Stores baseline prices for domestic accessorials for a GHC pricing contract (see sheet 5a).';
COMMENT ON COLUMN re_domestic_accessorial_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_domestic_accessorial_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_domestic_accessorial_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_domestic_accessorial_prices.services_schedule IS 'The services schedule (1, 2, or 3, based on location) for this price.';
COMMENT ON COLUMN re_domestic_accessorial_prices.per_unit_cents IS 'The price in cents, per unit of measure, for the service.';
COMMENT ON COLUMN re_domestic_accessorial_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_domestic_accessorial_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_domestic_linehaul_prices IS 'Stores baseline prices for domestic linehaul for a GHC pricing contract (see sheet 2a).';
COMMENT ON COLUMN re_domestic_linehaul_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_domestic_linehaul_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_domestic_linehaul_prices.weight_lower IS 'The lower bound of shipment weight (inclusive) for this price.';
COMMENT ON COLUMN re_domestic_linehaul_prices.weight_upper IS 'The upper bound of shipment weight (inclusive) for this price.';
COMMENT ON COLUMN re_domestic_linehaul_prices.miles_lower IS 'The lower bound of miles traveled (inclusive) for this price.';
COMMENT ON COLUMN re_domestic_linehaul_prices.miles_upper IS 'The upper bound of miles traveled (inclusive) for this price.';
COMMENT ON COLUMN re_domestic_linehaul_prices.is_peak_period IS 'Is this a peak period move?  Peak is May 15-Sept 30.';
COMMENT ON COLUMN re_domestic_linehaul_prices.domestic_service_area_id IS 'The domestic service area (based on zip3) for this price.';
COMMENT ON COLUMN re_domestic_linehaul_prices.price_millicents IS 'The price in millicents per hundred weight (CWT) per mile.';
COMMENT ON COLUMN re_domestic_linehaul_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_domestic_linehaul_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_domestic_other_prices IS 'Stores baseline prices for other domestic services for a GHC pricing contract (see sheet 2c).';
COMMENT ON COLUMN re_domestic_other_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_domestic_other_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_domestic_other_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_domestic_other_prices.is_peak_period IS 'Is this a peak period move?  Peak is May 15-Sept 30.';
COMMENT ON COLUMN re_domestic_other_prices.schedule IS 'The services schedule (1, 2, or 3, based on location) for this price.';
COMMENT ON COLUMN re_domestic_other_prices.price_cents IS 'The price in cents per hundred weight (CWT).';
COMMENT ON COLUMN re_domestic_other_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_domestic_other_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_domestic_service_area_prices IS 'Stores baseline prices for services within a domestic service area for a GHC pricing contract (see sheet 2b).';
COMMENT ON COLUMN re_domestic_service_area_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_domestic_service_area_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_domestic_service_area_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_domestic_service_area_prices.is_peak_period IS 'Is this a peak period move?  Peak is May 15-Sept 30.';
COMMENT ON COLUMN re_domestic_service_area_prices.domestic_service_area_id IS 'The domestic service area (based on zip3) for this price.';
COMMENT ON COLUMN re_domestic_service_area_prices.price_cents IS 'The price in cents. Some services are per hundred weight (CWT) per mile while others are just per hundred weight. See pricing template for details.';
COMMENT ON COLUMN re_domestic_service_area_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_domestic_service_area_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_domestic_service_areas IS 'Represents the domestic service areas defined in a GHC pricing contract (see sheet 1b).';
COMMENT ON COLUMN re_domestic_service_areas.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_domestic_service_areas.service_area IS 'A 3-digit code uniquely identifying a service area (e.g., 004, 344).';
COMMENT ON COLUMN re_domestic_service_areas.services_schedule IS 'The services schedule (1, 2, or 3) for this service area.';
COMMENT ON COLUMN re_domestic_service_areas.sit_pd_schedule IS 'The SIT (Storage In Transit) pickup/delivery schedule (1, 2, or 3) for this service area.';
COMMENT ON COLUMN re_domestic_service_areas.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_domestic_service_areas.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN re_domestic_service_areas.contract_id IS 'The associated GHC pricing contract.';

COMMENT ON TABLE re_intl_accessorial_prices IS 'Stores baseline prices for international accessorials for a GHC pricing contract (see sheet 5a).';
COMMENT ON COLUMN re_intl_accessorial_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_intl_accessorial_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_intl_accessorial_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_intl_accessorial_prices.market IS 'The market (CONUS or OCONUS) for this price.';
COMMENT ON COLUMN re_intl_accessorial_prices.per_unit_cents IS 'The price in cents, per unit of measure, for the service.';
COMMENT ON COLUMN re_intl_accessorial_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_intl_accessorial_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_intl_other_prices IS 'Stores baseline prices for other international services for a GHC pricing contract (see sheet 3d).';
COMMENT ON COLUMN re_intl_other_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_intl_other_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_intl_other_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_intl_other_prices.is_peak_period IS 'Is this a peak period move?  Peak is May 15-Sept 30.';
COMMENT ON COLUMN re_intl_other_prices.rate_area_id IS 'The rate area (based on location) for this price.';
COMMENT ON COLUMN re_intl_other_prices.per_unit_cents IS 'The price in cents. Some services are per hundred weight (CWT); others are per hundred weight per mile. See pricing template for details.';
COMMENT ON COLUMN re_intl_other_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_intl_other_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_intl_prices IS 'Stores baseline prices for international services for a GHC pricing contract (see sheets 3a-3c).';
COMMENT ON COLUMN re_intl_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_intl_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_intl_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_intl_prices.is_peak_period IS 'Is this a peak period move?  Peak is May 15-Sept 30.';
COMMENT ON COLUMN re_intl_prices.origin_rate_area_id IS 'The origin rate area (based on location) for this price.';
COMMENT ON COLUMN re_intl_prices.destination_rate_area_id IS 'The destination rate area (based on location) for this price.';
COMMENT ON COLUMN re_intl_prices.per_unit_cents IS 'The price in cents per hundred weight (CWT).';
COMMENT ON COLUMN re_intl_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_intl_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_rate_areas IS 'Represents the rate areas defined in a GHC pricing contract (see sheets 3a-3e).';
COMMENT ON COLUMN re_rate_areas.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_rate_areas.is_oconus IS 'Is this rate area for an OCONUS location?';
COMMENT ON COLUMN re_rate_areas.code IS 'A short alphanumeric code uniquely identifying a rate area (e.g., AR, GR29, US13).';
COMMENT ON COLUMN re_rate_areas.name IS 'A descriptive name for the rate area.';
COMMENT ON COLUMN re_rate_areas.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_rate_areas.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN re_rate_areas.contract_id IS 'The associated GHC pricing contract.';

COMMENT ON TABLE re_services IS 'Represents the move-related services that are included in a GHC pricing contract.';
COMMENT ON COLUMN re_services.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_services.code IS 'A short alphabetical code uniquely identifying a service (e.g., DLH, FSC)';
COMMENT ON COLUMN re_services.name IS 'A descriptive name for the service.';
COMMENT ON COLUMN re_services.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_services.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN re_services.priority IS 'The priority of this service in a payment request; a lower number indicates a higher priority (i.e., should be priced first).';

COMMENT ON TABLE re_shipment_type_prices IS 'Stores baseline prices for services associated with a shipment type for a GHC pricing contract (see sheet 5a).';
COMMENT ON COLUMN re_shipment_type_prices.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_shipment_type_prices.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_shipment_type_prices.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_shipment_type_prices.market IS 'The market (CONUS or OCONUS) for this price.';
COMMENT ON COLUMN re_shipment_type_prices.factor IS 'The price factor. Other domestic/international prices are multiplied by this factor. See pricing template for details.';
COMMENT ON COLUMN re_shipment_type_prices.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_shipment_type_prices.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_task_order_fees IS 'Stores prices for services associated with a task order for a GHC pricing contract (see sheet 4a).';
COMMENT ON COLUMN re_task_order_fees.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_task_order_fees.contract_year_id IS 'The associated contract year.';
COMMENT ON COLUMN re_task_order_fees.service_id IS 'The associated service being priced.';
COMMENT ON COLUMN re_task_order_fees.price_cents IS 'The price in cents per task order. Note that price escalations do not apply. See pricing template for details.';
COMMENT ON COLUMN re_task_order_fees.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_task_order_fees.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE re_zip3s IS 'Represents the zip3s defined in a GHC pricing contract (see sheet 1b) along with their associated service/rate areas.';
COMMENT ON COLUMN re_zip3s.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_zip3s.zip3 IS 'The first three digits of a zip code.';
COMMENT ON COLUMN re_zip3s.domestic_service_area_id IS 'The associated domestic service area for this zip3.';
COMMENT ON COLUMN re_zip3s.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_zip3s.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN re_zip3s.contract_id IS 'The associated GHC pricing contract.';
COMMENT ON COLUMN re_zip3s.rate_area_id IS 'The associated rate area for this zip3.';
COMMENT ON COLUMN re_zip3s.has_multiple_rate_areas IS 'True if this zip3 has multiple rate areas within it; if true, see the re_zip5_rate_areas table to determine the rate area.';
COMMENT ON COLUMN re_zip3s.base_point_city IS 'The name of the base point (primary) city associated with this zip3.';
COMMENT ON COLUMN re_zip3s.state IS 'The state for the base point city.';

COMMENT ON TABLE re_zip5_rate_areas IS 'Given a zip3 that has multiple rate areas, this table will associate the more-specific zip5 in that zip3 with a rate area.';
COMMENT ON COLUMN re_zip5_rate_areas.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN re_zip5_rate_areas.rate_area_id IS 'The associated rate area for this zip5.';
COMMENT ON COLUMN re_zip5_rate_areas.zip5 IS 'The full five-digit zip code.';
COMMENT ON COLUMN re_zip5_rate_areas.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN re_zip5_rate_areas.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN re_zip5_rate_areas.contract_id IS 'The associated GHC pricing contract.';

COMMENT ON TABLE schema_migration IS 'Stores the version (a date stamp in our case) for the database migrations that have been applied to this database.';
COMMENT ON COLUMN schema_migration.version IS 'A unique version string for the migration; derived from the first part of the migration filename.';

COMMENT ON TABLE service_item_param_keys IS 'Represents the keys for parameters that can be associated to a move-related service.';
COMMENT ON COLUMN service_item_param_keys.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN service_item_param_keys.key IS 'A short, human-readable string for the parameter.';
COMMENT ON COLUMN service_item_param_keys.description IS 'A descriptive name for the parameter.';
COMMENT ON COLUMN service_item_param_keys.type IS 'The type of the value associated with this key; options are STRING, DATE, INTEGER, DECIMAL, TIMESTAMP, PaymentServiceItemUUID.';
COMMENT ON COLUMN service_item_param_keys.origin IS 'Where values for this key originate; options are PRIME (the GHC prime contractor provides) or SYSTEM (the system determines the value).';
COMMENT ON COLUMN service_item_param_keys.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN service_item_param_keys.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE service_params IS 'Associates services with their expected input parameter keys.';
COMMENT ON COLUMN service_params.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN service_params.service_id IS 'The associated service.';
COMMENT ON COLUMN service_params.service_item_param_key_id IS 'The associated key.';
COMMENT ON COLUMN service_params.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN service_params.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_full_pack_rates IS 'Stores the rates for full pack from the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_full_pack_rates.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_full_pack_rates.schedule IS 'The services schedule (1, 2, 3, or 4, based on location) for this rate.';
COMMENT ON COLUMN tariff400ng_full_pack_rates.weight_lbs_lower IS 'The lower bound of shipment weight (inclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_full_pack_rates.weight_lbs_upper IS 'The upper bound of shipment weight (exclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_full_pack_rates.rate_cents IS 'The rate in cents per hundred weight (CWT).';
COMMENT ON COLUMN tariff400ng_full_pack_rates.effective_date_lower IS 'The start date for this rate (inclusive).';
COMMENT ON COLUMN tariff400ng_full_pack_rates.effective_date_upper IS 'The end date for this rate (exclusive).';
COMMENT ON COLUMN tariff400ng_full_pack_rates.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_full_pack_rates.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_full_unpack_rates IS 'Stores the rates for full unpack from the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.schedule IS 'The services schedule (1, 2, 3, or 4, based on location) for this rate.';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.rate_millicents IS 'The rate in millicents per hundred weight (CWT).';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.effective_date_lower IS 'The start date for this rate (inclusive).';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.effective_date_upper IS 'The end date for this rate (exclusive).';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_full_unpack_rates.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_item_rates IS 'Stores the rates for various item codes (accessorials) from the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_item_rates.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_item_rates.code IS 'The item code (e.g., 120C, 4B) for this rate.';
COMMENT ON COLUMN tariff400ng_item_rates.schedule IS 'The services schedule (1, 2, 3, or 4, based on location) for this rate, or null if rate is independent of schedule.';
COMMENT ON COLUMN tariff400ng_item_rates.weight_lbs_lower IS 'The lower bound of shipment weight (inclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_item_rates.weight_lbs_upper IS 'The upper bound of shipment weight (inclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_item_rates.rate_cents IS 'The rate in cents. The rates may be per hundred weight (CWT), per mile, per occurrence, or some combination. See 400NG tariff for details.';
COMMENT ON COLUMN tariff400ng_item_rates.effective_date_lower IS 'The start date for this rate (inclusive).';
COMMENT ON COLUMN tariff400ng_item_rates.effective_date_upper IS 'The end date for this rate (exclusive).';
COMMENT ON COLUMN tariff400ng_item_rates.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_item_rates.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_items IS 'Represents the items (accessorials) and their associated metadata from the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_items.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_items.code IS 'The item code (e.g., 120C, 4B)';
COMMENT ON COLUMN tariff400ng_items.discount_type IS 'The type of discount for this item; options are HHG, HHG_LINEHAUL_50, SIT, NONE.';
COMMENT ON COLUMN tariff400ng_items.allowed_location IS 'The allowed location for this item; options are ORIGIN, DESTINATION, NEITHER, EITHER.';
COMMENT ON COLUMN tariff400ng_items.item IS 'A descriptive name for this item.';
COMMENT ON COLUMN tariff400ng_items.measurement_unit_1 IS 'The first measurement unit needed for this item; options are BW (weight), CF (cubic foot), EA (each), FR (flat rate), FP (fuel percentage), NR (container), MV (monetary value), TD (days), TH (hours), NONE.';
COMMENT ON COLUMN tariff400ng_items.measurement_unit_2 IS 'The second measurement unit needed for this item; options are BW (weight), CF (cubic foot), EA (each), FR (flat rate), FP (fuel percentage), NR (container), MV (monetary value), TD (days), TH (hours), NONE.';
COMMENT ON COLUMN tariff400ng_items.rate_ref_code IS 'The reference code for this item; options are DD (date delivered), FS (fuel surcharge), MI (miles), PS (pack percentage), SC (point schedule), SE (tariff section), NONE.';
COMMENT ON COLUMN tariff400ng_items.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_items.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN tariff400ng_items.requires_pre_approval IS 'True if this item requires pre-approval.';

COMMENT ON TABLE tariff400ng_linehaul_rates IS 'Stores the rates for linehaul from the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.distance_miles_lower IS 'The lower bound of miles traveled (inclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.distance_miles_upper IS 'The upper bound of miles traveled (exclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.weight_lbs_lower IS 'The lower bound of shipment weight (inclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.weight_lbs_upper IS 'The upper bound of shipment weight (exclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.rate_cents IS 'The rate in cents.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.effective_date_lower IS 'The start date for this rate (inclusive).';
COMMENT ON COLUMN tariff400ng_linehaul_rates.effective_date_upper IS 'The end date for this rate (exclusive).';
COMMENT ON COLUMN tariff400ng_linehaul_rates.type IS 'The type of this rate; options are ConusLinehaul and IntraAlaskaLinehaul.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_linehaul_rates.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_service_areas IS 'Represents the service areas defined in the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_service_areas.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_service_areas.service_area IS 'A code uniquely identifying a service area (e.g., 4, 344).';
COMMENT ON COLUMN tariff400ng_service_areas.name IS 'The primary city/state for this service area.';
COMMENT ON COLUMN tariff400ng_service_areas.services_schedule IS 'The services schedule (1, 2, 3, or 4) for this service area.';
COMMENT ON COLUMN tariff400ng_service_areas.linehaul_factor IS 'The factor in cents per hundred weight (CWT) to apply to the linehaul rate.';
COMMENT ON COLUMN tariff400ng_service_areas.service_charge_cents IS 'The 135A/135B origin/destination service charge in cents per hundred weight.';
COMMENT ON COLUMN tariff400ng_service_areas.effective_date_lower IS 'The start date for this rate (inclusive).';
COMMENT ON COLUMN tariff400ng_service_areas.effective_date_upper IS 'The end date for this rate (exclusive).';
COMMENT ON COLUMN tariff400ng_service_areas.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_service_areas.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN tariff400ng_service_areas.sit_185a_rate_cents IS 'The rate in cents for the SIT first day (185A) item.';
COMMENT ON COLUMN tariff400ng_service_areas.sit_185b_rate_cents IS 'The rate in cents for the SIT additional days (185B) item.';
COMMENT ON COLUMN tariff400ng_service_areas.sit_pd_schedule IS 'The SIT pickup/delivery schedule (1, 2, 3, or 4) for this service area.';

COMMENT ON TABLE tariff400ng_shorthaul_rates IS 'Stores the rates for shorthaul from the 400NG tariff.';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.cwt_miles_lower IS 'The lower bound of hundred weight (CWT) multiplied by miles traveled (inclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.cwt_miles_upper IS 'The upper bound of hundred weight (CWT) multiplied by miles traveled (exclusive) for this rate.';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.rate_cents IS 'The rate in cents.';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.effective_date_lower IS 'The start date for this rate (inclusive).';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.effective_date_upper IS 'The end date for this rate (exclusive).';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_shorthaul_rates.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_zip3s IS 'Represents the zip3s defined in the 400NG tariff along with their associated service areas, rate areas, and regions.';
COMMENT ON COLUMN tariff400ng_zip3s.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_zip3s.zip3 IS 'The first three digits of a zip code.';
COMMENT ON COLUMN tariff400ng_zip3s.basepoint_city IS 'The name of the base point (primary) city associated with this zip3.';
COMMENT ON COLUMN tariff400ng_zip3s.state IS 'The state for the base point city.';
COMMENT ON COLUMN tariff400ng_zip3s.service_area IS 'The associated service area (e.g., 56, 184) for this zip3.';
COMMENT ON COLUMN tariff400ng_zip3s.rate_area IS 'The associated rate area (e.g., US20, US47) for this zip3. If the rate area is ZIP, then the zip3 is not sufficient to determine the rate area. In that case, use the zip5 along with the tariff400ng_zip5_rate_areas table to determine the rate area.';
COMMENT ON COLUMN tariff400ng_zip3s.region IS 'The associated region (e.g., 13, 6) for this zip3.';
COMMENT ON COLUMN tariff400ng_zip3s.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_zip3s.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE tariff400ng_zip5_rate_areas IS 'Given a zip3 that has multiple rate areas, this table will associate the more-specific zip5 in that zip3 with a rate area.';
COMMENT ON COLUMN tariff400ng_zip5_rate_areas.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN tariff400ng_zip5_rate_areas.zip5 IS 'The full five-digit zip code.';
COMMENT ON COLUMN tariff400ng_zip5_rate_areas.rate_area IS 'The associated rate area for this zip5.';
COMMENT ON COLUMN tariff400ng_zip5_rate_areas.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN tariff400ng_zip5_rate_areas.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE traffic_distribution_lists IS 'Represents the possible channels (rate area to region for a code of service) for a move.';
COMMENT ON COLUMN traffic_distribution_lists.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN traffic_distribution_lists.source_rate_area IS 'The rate area for the origin.';
COMMENT ON COLUMN traffic_distribution_lists.destination_region IS 'The region for the destination.';
COMMENT ON COLUMN traffic_distribution_lists.code_of_service IS 'The code of service for this channel; options are D and 2. See 400NG tariff for details.';
COMMENT ON COLUMN traffic_distribution_lists.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN traffic_distribution_lists.updated_at IS 'Timestamp when the record was last updated.';

COMMENT ON TABLE transportation_service_provider_performances IS 'Stores scores/rates for transportation service providers (TSPs) on a given channel.';
COMMENT ON COLUMN transportation_service_provider_performances.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN transportation_service_provider_performances.performance_period_start IS 'The start date of the performance period (inclusive).';
COMMENT ON COLUMN transportation_service_provider_performances.performance_period_end IS 'The end date of the performance period (inclusive).';
COMMENT ON COLUMN transportation_service_provider_performances.traffic_distribution_list_id IS 'The associated traffic distribution list (or channel) for this performance.';
COMMENT ON COLUMN transportation_service_provider_performances.quality_band IS 'The quality band (1 through 4) that this performance falls within; previously used by the HHG award queue.';
COMMENT ON COLUMN transportation_service_provider_performances.offer_count IS 'How many offers have been made for this performance; previously used by the HHG award queue.';
COMMENT ON COLUMN transportation_service_provider_performances.best_value_score IS 'The best value score (BVS) for this performance; a ranking that was previously used by the HHG award queue.';
COMMENT ON COLUMN transportation_service_provider_performances.transportation_service_provider_id IS 'The associated provider (moving company) for this performance.';
COMMENT ON COLUMN transportation_service_provider_performances.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN transportation_service_provider_performances.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN transportation_service_provider_performances.rate_cycle_start IS 'The start date of the rate cycle (inclusive).';
COMMENT ON COLUMN transportation_service_provider_performances.rate_cycle_end IS 'The end date of the rate cycle (inclusive).';
COMMENT ON COLUMN transportation_service_provider_performances.linehaul_rate IS 'The discount rate (taken off of the regular rate) for linehaul offered by the provider of this performance.';
COMMENT ON COLUMN transportation_service_provider_performances.sit_rate IS 'The discount rate (taken off of the regular rate) for storage-in-transit (SIT) offered by the provider of this performance.';

COMMENT ON TABLE transportation_service_providers IS 'Represents the transportation service providers (TSPs, or moving companies) used in pre-GHC HHG moves.';
COMMENT ON COLUMN transportation_service_providers.id IS 'UUID that uniquely identifies the record.';
COMMENT ON COLUMN transportation_service_providers.standard_carrier_alpha_code IS 'A unique two-to-four letter identifier (usually four-letter) for the TSP.';
COMMENT ON COLUMN transportation_service_providers.created_at IS 'Timestamp when the record was first created.';
COMMENT ON COLUMN transportation_service_providers.updated_at IS 'Timestamp when the record was last updated.';
COMMENT ON COLUMN transportation_service_providers.enrolled IS 'True if this TSP is able to accept offers for HHG moves.';
COMMENT ON COLUMN transportation_service_providers.name IS 'A descriptive name for this TSP.';
COMMENT ON COLUMN transportation_service_providers.supplier_id IS 'The supplier ID for this TSP.';
COMMENT ON COLUMN transportation_service_providers.poc_general_name IS 'A point-of-contact name for general inquiries.';
COMMENT ON COLUMN transportation_service_providers.poc_general_email IS 'A point-of-contact email for general inquiries.';
COMMENT ON COLUMN transportation_service_providers.poc_general_phone IS 'A point-of-contact phone number for general inquiries.';
COMMENT ON COLUMN transportation_service_providers.poc_claims_name IS 'A point-of-contact name for claims inquiries.';
COMMENT ON COLUMN transportation_service_providers.poc_claims_email IS 'A point-of-contact email for claims inquiries.';
COMMENT ON COLUMN transportation_service_providers.poc_claims_phone IS 'A point-of-contact phone number for claims inquiries.';
