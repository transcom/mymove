ALTER TABLE addresses
	ADD CONSTRAINT us_post_region_cities_id_fkey
	FOREIGN KEY (us_post_region_cities_id) REFERENCES us_post_region_cities (id);

DROP FUNCTION IF EXISTS calculate_escalated_price(uuid, uuid, uuid, uuid);

-- function to get the rate area id
CREATE OR REPLACE FUNCTION get_rate_area_id(
    address_id UUID,
    service_item_id UUID,
    OUT o_rate_area_id UUID
)
RETURNS UUID AS $$
BEGIN
    SELECT ro.rate_area_id
    INTO o_rate_area_id
    FROM addresses a
    JOIN re_oconus_rate_areas ro
    ON a.us_post_region_cities_id = ro.us_post_region_cities_id
    WHERE a.id = address_id;

    IF o_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate area not found for address % for service item id %', address_id, service_item_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- function to get the contract id based off of a specific date that falls between start/end
CREATE OR REPLACE FUNCTION get_contract_id(
    requested_pickup_date DATE,
    OUT o_contract_id UUID
)
RETURNS UUID AS $$
BEGIN
    -- Get the contract_id from the re_contract_years table
    SELECT rcy.contract_id
    INTO o_contract_id
    FROM re_contract_years rcy
    WHERE requested_pickup_date BETWEEN rcy.start_date AND rcy.end_date;

    -- Check if contract_id is found, else raise an exception
    IF o_contract_id IS NULL THEN
        RAISE EXCEPTION 'Contract not found for requested pickup date %', requested_pickup_date;
    END IF;
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION calculate_escalated_price(
    o_rate_area_id UUID,
    d_rate_area_id UUID,
    re_service_id UUID,
    c_id UUID
) RETURNS NUMERIC AS $$
DECLARE
    per_unit_cents NUMERIC;
    escalation_factor NUMERIC;
    escalated_price NUMERIC;
BEGIN
    SELECT rip.per_unit_cents
    INTO per_unit_cents
    FROM re_intl_prices rip
    WHERE rip.origin_rate_area_id = o_rate_area_id
      AND rip.destination_rate_area_id = d_rate_area_id
      AND rip.service_id = re_service_id
      AND rip.contract_id = c_id;

    -- IF per_unit_cents IS NULL THEN
    --     RAISE EXCEPTION 'No matching price found for the given parameters';
    -- END IF;

    SELECT rcy.escalation
    INTO escalation_factor
    FROM re_contract_years rcy
    WHERE rcy.contract_id = c_id;

    IF escalation_factor IS NULL THEN
        RAISE EXCEPTION 'Escalation factor not found for contract_id %', contract_id;
    END IF;

    -- Calculate the escalated price
    escalated_price := ROUND(per_unit_cents * escalation_factor::NUMERIC / 100, 2);

    RETURN escalated_price;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE PROCEDURE update_service_item_pricing(shipment_id UUID) AS
'
DECLARE
    shipment RECORD;
    service_item RECORD;
    escalated_price NUMERIC;
    estimated_price NUMERIC;
    o_rate_area_id UUID;
    d_rate_area_id UUID;
    contract_id UUID;
BEGIN
    SELECT ms.id, ms.pickup_address_id, ms.destination_address_id, ms.requested_pickup_date, ms.prime_estimated_weight
    INTO shipment
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF shipment IS NULL THEN
        RAISE EXCEPTION ''Shipment with ID % not found'', shipment_id;
    END IF;

    -- Loop through service items in the shipment and update pricing
    FOR service_item IN
        SELECT si.id, si.re_service_id
        FROM mto_service_items si
        WHERE si.mto_shipment_id = shipment_id
    LOOP
        -- Get origin and destination rate areas
        o_rate_area_id := get_rate_area_id(shipment.pickup_address_id, service_item.re_service_id);
        d_rate_area_id := get_rate_area_id(shipment.destination_address_id, service_item.re_service_id);
        contract_id := get_contract_id(shipment.requested_pickup_date);

        -- Calculate the escalated price
        escalated_price := calculate_escalated_price(
            o_rate_area_id,
            d_rate_area_id,
            service_item.re_service_id,
            contract_id
        );

        -- Calculate estimated and actual prices
        IF shipment.prime_estimated_weight IS NOT NULL THEN
            estimated_price := ROUND(escalated_price * (shipment.prime_estimated_weight / 100)::NUMERIC, 2);
        END IF;

        -- Update the pricing_estimate in mto_service_items
        UPDATE mto_service_items
        SET pricing_estimate = estimated_price
        WHERE id = service_item.id;
    END LOOP;
END;
'
LANGUAGE plpgsql;
