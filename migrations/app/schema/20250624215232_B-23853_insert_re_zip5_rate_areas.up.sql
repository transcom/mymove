--insert missing zip5 rate area association for Titusville, FL

insert into re_zip5_rate_areas (id, rate_area_id, zip5, created_at, updated_at, contract_id)
values ('1f272c56-7a97-4bae-9a9f-2a6f46047bea', 'ddd74fb8-c0f1-41a9-9d4f-234bd295ae1a', '32780', now(), now(), '070f7c82-fad0-4ae8-9a83-5de87a56472e')
on conflict (zip5, rate_area_id) do nothing;

insert into re_zip5_rate_areas (id, rate_area_id, zip5, created_at, updated_at, contract_id)
values ('c47bc035-77fb-448e-9d5e-2512f2904ccd', 'ddd74fb8-c0f1-41a9-9d4f-234bd295ae1a', '32796', now(), now(), '070f7c82-fad0-4ae8-9a83-5de87a56472e')
on conflict (zip5, rate_area_id) do nothing;