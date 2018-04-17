ALTER TABLE tariff400ng_service_areas ADD COLUMN sit_185A_rate_cents int;
ALTER TABLE tariff400ng_service_areas ADD COLUMN sit_185B_rate_cents int;
ALTER TABLE tariff400ng_service_areas ADD COLUMN sit_pd_schedule int;

CREATE FUNCTION import_sit_rates() RETURNS void AS $$
DECLARE
	temp_sit_rate record;
BEGIN
	FOR temp_sit_rate IN SELECT * FROM temp_sit_rates LOOP
		UPDATE tariff400ng_service_areas
		SET
			sit_185A_rate_cents = temp_sit_rate.sit_185A_rate_cents,
			sit_185B_rate_cents = temp_sit_rate.sit_185a_rate_cents,
			sit_pd_schedule = temp_sit_rate.sit_pd_schedule
		WHERE tariff400ng_service_areas.service_area = temp_sit_rate.service_area_number
			AND tariff400ng_service_areas.effective_date_lower = '2018-05-15';
	END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT import_sit_rates();

DROP TABLE temp_sit_rates;
DROP FUNCTION import_sit_rates;
