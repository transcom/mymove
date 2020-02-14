create table re_zip5_rate_areas
(
    id uuid
        constraint re_zip5_rate_areas_pkey primary key,
    rate_area_id uuid not null
        constraint re_zip5_rate_areas_rate_area_id_fkey references re_rate_areas,
    zip5 text not null
        constraint re_zip5_rate_areas_zip5_key unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

alter table re_zip3s
    add rate_area_id uuid
    constraint re_zip3s_rate_area_id_fkey references re_rate_areas,
    add has_multiple_rate_areas boolean not null default false
    constraint re_zip3s_has_zip5 CHECK (
        (has_multiple_rate_areas = true and rate_area_id is null) or (has_multiple_rate_areas = false and rate_area_id is not null)
    );