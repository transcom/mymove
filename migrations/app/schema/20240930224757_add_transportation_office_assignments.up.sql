CREATE TABLE IF NOT EXISTS public.transportation_office_assignments (
    id uuid NOT NULL,
	transportation_office_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT null,
    primary_office bool
);

INSERT INTO public.transportation_office_assignments SELECT id, transportation_office_id, created_at, updated_at FROM office_users;

UPDATE public.transportation_office_assignments toa SET primary_office = true FROM public.office_users ofusr WHERE ofusr.transportation_office_id = toa.transportation_office_id;
