ALTER TABLE public.mto_shipments
ADD COLUMN IF NOT EXISTS prime_acknowledged_at TIMESTAMPTZ NULL;

COMMENT ON COLUMN public.mto_shipments.prime_acknowledged_at IS 'The date and time the prime acknowledged the shipment';