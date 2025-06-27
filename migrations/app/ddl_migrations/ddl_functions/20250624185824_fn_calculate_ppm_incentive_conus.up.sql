-- B-23853 Beth Grohmann Initial check-in

DROP FUNCTION IF EXISTS public.calculate_ppm_incentive_conus(uuid, bool, bool, bool, bool);

CREATE OR REPLACE FUNCTION public.calculate_ppm_incentive_conus(ppm_id uuid, is_estimated boolean, is_actual boolean, is_max boolean, update_table boolean)
  RETURNS TABLE(total_incentive numeric, price_dsh numeric, price_dlh numeric, price_ddp numeric, price_dop numeric, price_dpk numeric, price_dupk numeric, price_fsc numeric)
  LANGUAGE plpgsql
AS $function$
declare
  ppm RECORD;
  v_contract_id UUID;
  o_rate_area_id UUID;
  d_rate_area_id UUID;
  service_id UUID;
  estimated_fsc_multiplier numeric;
  fuel_price numeric;
  price_difference numeric;
  cents_above_baseline numeric;
  peak_period BOOLEAN;
  pickup_address_id UUID;
  destination_address_id UUID;
  mileage integer;
  weight integer;
  move_date date;
  raw_millicents numeric;
  weight_lower_val numeric;
  weight_upper_val numeric;
  miles_lower_val numeric;
  miles_upper_val numeric;
  cents_per_cwt numeric;
  escalation_factor numeric;
  gcc_multiplier NUMERIC := 1.00;
  v_gcc_multiplier_id uuid;
  grade text;
  pickup_zip3 text;
  destination_zip3 text;
begin

  if not is_estimated
  and not is_actual
  and not is_max then
    raise exception 'is_estimated, is_actual, and is_max cannot all be FALSE. No update will be performed.';
  end if;

  select
    ppms.id,
    ppms.max_incentive
  into
    ppm
  from
    ppm_shipments ppms
  where ppms.id = ppm_id;

  if ppm is null then
    raise exception 'PPM with ID % not found', ppm_id;
  end if;

  SELECT
    ppm_shipments.pickup_postal_address_id,
    ppm_shipments.destination_postal_address_id,
    coalesce(ppm_shipments.actual_move_date,ppm_shipments.expected_departure_date),
    mto_shipments.distance,
    ppm_shipments.estimated_weight,
    f.grade
  INTO
    pickup_address_id,
    destination_address_id,
    move_date,
    mileage,
    weight,
    grade
  FROM
    ppm_shipments
  LEFT JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
  LEFT JOIN moves d on mto_shipments.move_id = d.id
  LEFT JOIN orders f on d.orders_id = f.id
  LEFT JOIN service_members e on f.service_member_id = e.id
  WHERE ppm_shipments.id = ppm_id;

  IF is_actual THEN
    SELECT
      COALESCE(SUM(
        CASE
          WHEN adjusted_net_weight IS NOT NULL THEN adjusted_net_weight
          ELSE full_weight - empty_weight
        END), 0)
    INTO weight
    FROM weight_tickets wt
    WHERE wt.ppm_shipment_id = ppm_id
      AND (wt.status NOT IN ('REJECTED', 'EXCLUDED') OR wt.status IS NULL);
  END IF;

  peak_period := is_peak_period(move_date);

  v_contract_id := get_contract_id(move_date);
  if v_contract_id is null then
    raise exception 'Contract not found for date: %',move_date;
  end if;

  o_rate_area_id := get_rate_area_id(pickup_address_id,null,v_contract_id);
  if o_rate_area_id is null then
    raise exception 'Origin rate area is NULL for address ID %',pickup_address_id;
  end if;

  d_rate_area_id := get_rate_area_id(destination_address_id,null,v_contract_id);
  if d_rate_area_id is null then
    raise exception 'Destination rate area is NULL for address ID %',destination_address_id;
  end if;

  -- get the ZIP3s of pickup and destination addresses
  -- determines if we cost shorthaul or linehaul
  SELECT LEFT(a.postal_code, 3) INTO pickup_zip3
  FROM addresses a
  WHERE a.id = pickup_address_id;

  SELECT LEFT(a.postal_code, 3) INTO destination_zip3
  FROM addresses a
  WHERE a.id = destination_address_id;

  SELECT escalation_compounded INTO escalation_factor
    FROM re_contract_years AS rcy
  WHERE rcy.contract_id = v_contract_id
    AND move_date BETWEEN rcy.start_date AND rcy.end_date;

  IF pickup_zip3 = destination_zip3 THEN
    price_dlh := 0;
    -- calculate DSH if ZIP3s are the same
    service_id := get_service_id('DSH');
    SELECT price_cents
    INTO price_dsh
    FROM re_domestic_service_area_prices dsap
    JOIN re_domestic_service_areas sa ON dsap.domestic_service_area_id = sa.id
    JOIN re_contracts c ON c.id = dsap.contract_id
    JOIN re_zip3s rz ON pickup_zip3 = rz.zip3
    WHERE dsap.contract_id = v_contract_id
    AND dsap.service_id = (SELECT id FROM re_services WHERE code = 'DSH')
    AND dsap.is_peak_period = peak_period
    AND dsap.domestic_service_area_id = sa.id
    LIMIT 1;

    price_dsh := ROUND(price_dsh * escalation_factor, 3);
    --RAISE NOTICE 'DSH price with escalation factor: %', price_dsh;

    price_dsh := ROUND(price_dsh * (weight::NUMERIC / 100) * mileage, 0);
    --RAISE NOTICE 'DSH final price: %', price_dsh;
  ELSE
    price_dsh := 0;
    -- calculate DLH instead
    service_id := get_service_id('DLH');
    SELECT rdlp.price_millicents
      INTO raw_millicents
      FROM re_domestic_linehaul_prices AS rdlp
      JOIN re_domestic_service_areas AS sa on rdlp.domestic_service_area_id = sa.id
      JOIN re_zip3s AS rzs ON sa.id = rzs.domestic_service_area_id
      JOIN addresses AS a ON LEFT(a.postal_code, 3) = rzs.zip3
     WHERE rdlp.contract_id = v_contract_id
       AND rdlp.is_peak_period = peak_period
       AND weight BETWEEN rdlp.weight_lower AND rdlp.weight_upper
       AND mileage BETWEEN rdlp.miles_lower AND rdlp.miles_upper
       AND a.id = pickup_address_id;

    --RAISE NOTICE 'DLH raw_millicents: %', raw_millicents;
    cents_per_cwt := ROUND(raw_millicents / 1000.0, 1);
    --RAISE NOTICE 'DLH cents_per_cwt: %', cents_per_cwt;

    cents_per_cwt := ROUND(cents_per_cwt * escalation_factor, 3);
    --RAISE NOTICE 'DLH cents_per_cwt with escalation factor: %', cents_per_cwt;

    price_dlh := ROUND(cents_per_cwt * (weight::NUMERIC / 100) * mileage, 0);
    --RAISE NOTICE 'DLH final price: %', price_dlh;
  END IF;

  -- Calculate DOP price
  service_id := get_service_id('DOP');
  price_dop := calculate_escalated_price_domestic(
        o_rate_area_id,
        null,
        service_id,
        v_contract_id,
        'DOP',
        move_date,
        pickup_address_id
      );
  --RAISE NOTICE 'DOP price (before weight): %', price_dop;
  price_dop := price_dop * (weight::numeric / 100);
  --RAISE NOTICE 'DOP price (after weight): %', price_dop;
  price_dop := ROUND(price_dop * 100);
  --RAISE NOTICE 'DOP price (after * 100): %', price_dop;

  service_id := get_service_id('DUPK');
  price_dupk := calculate_escalated_price_domestic(
      null,
      d_rate_area_id,
      service_id,
      v_contract_id,
      'DUPK',
      move_date,
      destination_address_id
    );
  --RAISE NOTICE 'DUPK price (before weight): %', price_dupk;
  price_dupk := price_dupk * (weight::numeric / 100);
  --RAISE NOTICE 'DUPK price (after weight): %', price_dupk;
  price_dupk := ROUND(price_dupk * 100);
  --RAISE NOTICE 'DUPK price (after * 100): %', price_dupk;

  -- Calculate DPK price
  service_id := get_service_id('DPK');
  price_dpk := calculate_escalated_price_domestic(
        o_rate_area_id,
        null,
        service_id,
        v_contract_id,
        'DPK',
        move_date,
        pickup_address_id
      );
  --RAISE NOTICE 'DPK price (before weight): %', price_dpk;
  price_dpk := price_dpk * (weight::numeric / 100);
  --RAISE NOTICE 'DPK price (after weight): %', price_dpk;
  price_dpk := ROUND(price_dpk * 100);
  --RAISE NOTICE 'DPK price (after * 100): %', price_dpk;

  -- Calculate DDP price
  service_id := get_service_id('DDP');
  price_ddp := calculate_escalated_price_domestic(
        null,
        d_rate_area_id,
        service_id,
        v_contract_id,
        'DDP',
        move_date,
        destination_address_id
      );
  --RAISE NOTICE 'DDP price (before weight): %', price_ddp;
  price_ddp := price_ddp * (weight::numeric / 100);
  --RAISE NOTICE 'DDP price (after weight): %', price_ddp;
  price_ddp := ROUND(price_ddp * 100);
  --RAISE NOTICE 'DDP price (after * 100): %', price_ddp;

  -- Calculate FSC price
  estimated_fsc_multiplier := get_fsc_multiplier(weight);
  fuel_price := get_fuel_price(move_date);
  price_difference := calculate_price_difference(fuel_price);
  cents_above_baseline := mileage * estimated_fsc_multiplier;
  price_fsc := ROUND((cents_above_baseline * price_difference) * 100);

  IF grade != 'CIVILIAN_EMPLOYEE' THEN --do not apply multiplier for civilians

    EXECUTE 'SELECT multiplier, id FROM gcc_multipliers WHERE $1 BETWEEN start_date AND end_date LIMIT 1' INTO gcc_multiplier, v_gcc_multiplier_id USING move_date;
    --RAISE NOTICE 'GCC Multiplier %', gcc_multiplier;

    IF price_dsh > 0 AND gcc_multiplier != 1.00 THEN
      price_dsh := ROUND(price_dsh * gcc_multiplier);
    END IF;

    IF price_dlh > 0 AND gcc_multiplier != 1.00 THEN
      price_dlh := ROUND(price_dlh * gcc_multiplier);
    END IF;

    IF price_dop > 0 AND gcc_multiplier != 1.00 THEN
      price_dop := ROUND(price_dop * gcc_multiplier);
    END IF;
    --raise notice 'DOP price after multiplier: %', price_dop;

    IF price_ddp > 0 AND gcc_multiplier != 1.00 THEN
      price_ddp := ROUND(price_ddp * gcc_multiplier);
    END IF;

    IF price_dpk > 0 AND gcc_multiplier != 1.00 THEN
      price_dpk := ROUND(price_dpk * gcc_multiplier);
    END IF;

    IF price_dupk > 0 AND gcc_multiplier != 1.00 THEN
      price_dupk := ROUND(price_dupk * gcc_multiplier);
    END IF;

    IF price_fsc > 0 AND gcc_multiplier != 1.00 THEN
      price_fsc := ROUND(price_fsc * gcc_multiplier);
    END IF;

  END IF;

  -- Calculate total incentive
  total_incentive :=
      COALESCE(price_dsh, 0) +
      COALESCE(price_dlh, 0) +
      COALESCE(price_dop, 0) +
      COALESCE(price_ddp, 0) +
      COALESCE(price_dpk, 0) +
      COALESCE(price_dupk, 0) +
      COALESCE(price_fsc, 0);

  -- we want to cap the final incentive to not be greater than the max incentive
  IF total_incentive > COALESCE(ppm.max_incentive, 0) AND is_actual THEN
    total_incentive := COALESCE(ppm.max_incentive, 0);
  END IF;

  IF update_table THEN
    UPDATE ppm_shipments
    SET estimated_incentive = CASE WHEN is_estimated THEN total_incentive ELSE estimated_incentive END,
    final_incentive = CASE WHEN is_actual THEN total_incentive ELSE final_incentive END,
    max_incentive = CASE WHEN is_max THEN total_incentive ELSE max_incentive END,
    gcc_multiplier_id = v_gcc_multiplier_id
    WHERE id = ppm_id;
  END IF;

  return QUERY
  select
    total_incentive,
    price_dsh,
    price_dlh,
    price_ddp,
    price_dop,
    price_dpk,
    price_dupk,
    price_fsc;

END;
$function$
;