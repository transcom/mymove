-- B-23853 Beth Grohmann Initial check-in

DROP FUNCTION IF EXISTS public.calculate_escalated_price_domestic(uuid, uuid, uuid, uuid, text, date, uuid);

CREATE OR REPLACE FUNCTION public.calculate_escalated_price_domestic(o_rate_area_id uuid, d_rate_area_id uuid, re_service_id uuid, c_id uuid, service_code text, requested_pickup_date date, address_id uuid)
 RETURNS numeric
 LANGUAGE plpgsql
AS $function$
declare
    per_unit_cents numeric;
    escalation_factor numeric;
    escalated_price numeric;
    is_oconus BOOLEAN;
    peak_period BOOLEAN;
begin
    peak_period := is_peak_period(requested_pickup_date);
    --raise notice 'ORIGIN: %', o_rate_area_id;
    --raise notice 'domestic: %', d_rate_area_id;
    --raise notice 'service_code: %', service_code;
    --raise notice 'pick up date: %', requested_pickup_date;
    --raise notice 'Contract: %', c_id;
    --raise notice 'SERIVCE ID: %', re_service_id;
    if service_code in ('DPK', 'DUPK') then
        select rip.price_cents
        into per_unit_cents
        from re_domestic_other_prices rip
        join re_domestic_service_areas sa  on sa.services_schedule = rip.schedule
        join re_zip3s rzs on sa.id = rzs.domestic_service_area_id
        join addresses a on left(a.postal_code, 3) = rzs.zip3
        where  rip.service_id = re_service_id
        and rip.contract_id = c_id
        and rip.is_peak_period = peak_period
        and a.id = address_id;
    else
        select dsap.price_cents
        into per_unit_cents
        from re_domestic_service_area_prices dsap
        join re_domestic_service_areas sa on dsap.domestic_service_area_id = sa.id
        join re_zip3s rzs on sa.id = rzs.domestic_service_area_id
        join addresses a on left(a.postal_code, 3) = rzs.zip3
        where dsap.service_id = re_service_id
        and dsap.contract_id = c_id
        and dsap.is_peak_period = peak_period
        and a.id = address_id;
    end if;
    --raise notice '% per unit cents: %', service_code, per_unit_cents;
    if per_unit_cents is null then
        raise exception 'No per unit cents found for service item id: %, origin rate area: %, dest rate area: %, and contract_id: %', 
            re_service_id, o_rate_area_id, d_rate_area_id, c_id;
    end if;
    select rcy.escalation_compounded
    into escalation_factor
    from re_contract_years rcy
    where rcy.contract_id = c_id
    and requested_pickup_date between rcy.start_date and rcy.end_date;
    if escalation_factor is null then
        raise exception 'Escalation factor not found for contract_id %', c_id;
    end if;
    per_unit_cents := per_unit_cents::numeric / 100.0;
    escalated_price := ROUND(per_unit_cents * escalation_factor, 2);
    return escalated_price;
end $function$;