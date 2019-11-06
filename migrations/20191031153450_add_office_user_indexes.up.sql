CREATE INDEX office_users_first_name_idx ON public.office_users USING btree (first_name);
CREATE INDEX office_users_last_name_idx ON public.office_users USING btree (last_name);
CREATE INDEX office_users_email_trgrm_idx ON public.office_users USING gin(email gin_trgm_ops);
