-- Make standard_carrier_alpha_code on transportation_service_providers unique --

ALTER TABLE transportation_service_providers
  ADD CONSTRAINT unique_standard_carrier_alpha_code UNIQUE (standard_carrier_alpha_code);
