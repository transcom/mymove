ALTER TABLE public.ppm_shipments
ADD COLUMN pickup_postal_address_id uuid NULL,
ADD COLUMN secondary_pickup_postal_address_id uuid NULL,
ADD COLUMN destination_postal_address_id uuid NULL,
ADD COLUMN secondary_destination_postal_address_id uuid NULL;
ALTER TABLE public.ppm_shipments
ADD CONSTRAINT ppm_shipments_pickup_postal_address_fkey FOREIGN KEY (pickup_postal_address_id) REFERENCES public.addresses(id),
ADD CONSTRAINT ppm_shipments_secondary_pickup_postal_address_fkey FOREIGN KEY (secondary_pickup_postal_address_id) REFERENCES public.addresses(id),
ADD CONSTRAINT ppm_shipments_destination_postal_address_fkey FOREIGN KEY (destination_postal_address_id) REFERENCES public.addresses(id),
ADD CONSTRAINT ppm_shipments_secondary_destination_postal_address_fkey FOREIGN KEY (secondary_destination_postal_address_id) REFERENCES public.addresses(id);