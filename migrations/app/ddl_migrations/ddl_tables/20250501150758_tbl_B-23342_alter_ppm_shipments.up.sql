alter table ppm_shipments
  add column if not exists has_gun_safe bool,
  add column if not exists gun_safe_weight int4;

comment on column ppm_shipments.has_gun_safe is 'Flag to indicate if PPM shipment has a gun safe';
comment on column ppm_shipments.gun_safe_weight is 'Customer estimated gun safe weight';