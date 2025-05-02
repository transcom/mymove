-- B-22831  Beth Grohmann  Add calculate_ppm_closeout
-- B-23586  Alex Lusk  Modify calculate_ppm_closeout to run as a migration, add skip_if_existsm, add zero-fill
DROP PROCEDURE IF EXISTS calculate_ppm_closeout;

CREATE OR REPLACE PROCEDURE calculate_ppm_closeout(p_ppm_shipment_id UUID, skip_if_exists bool)
LANGUAGE plpgsql
AS '
DECLARE
	i								record;
	v_ppm_closeout_id				uuid;
	v_estimated_incentive			int4;
	v_advance_amount_received		int4;
	v_final_incentive				int4;
	v_max_incentive 				int4;
	v_ppm_type						text;
	v_max_advance					int4;
	v_remaining_incentive			int4;
	v_total_gtcc_paid_expenses		int4;
	v_total_member_paid_expenses	int4;
	v_gtcc_paid_sit					int4;
	v_member_paid_sit				int4;
	v_gtcc_disbursement				int4;
	v_member_disbursement			int4;
	v_sql							text;
	v_count							int;
BEGIN
	SELECT COUNT(*) INTO v_count
	FROM ppm_closeouts
	WHERE ppm_shipment_id = p_ppm_shipment_id;

	IF v_count = 0 THEN
		INSERT INTO ppm_closeouts (id, ppm_shipment_id, created_at, updated_at)
		VALUES (uuid_generate_v4(), p_ppm_shipment_id, now(), now())
		RETURNING id INTO v_ppm_closeout_id;
	ELSE
		IF skip_if_exists THEN
			RETURN;
		ELSE
			SELECT id INTO v_ppm_closeout_id
			FROM ppm_closeouts
			WHERE ppm_shipment_id = p_ppm_shipment_id;
		END IF;
	END IF;

	RAISE NOTICE ''ppm_closeout_id: %'', v_ppm_closeout_id;

	--get gtcc paid expenses
	FOR i IN (
		SELECT moving_expense_type, SUM(amount) sum_amt
		FROM moving_expenses
		WHERE ppm_shipment_id = p_ppm_shipment_id
		   	AND status = ''APPROVED''
		   	AND paid_with_gtcc = true
		   	AND moving_expense_type != ''STORAGE''
		GROUP BY moving_expense_type)
	LOOP
		v_sql := format(
			''UPDATE ppm_closeouts
			 SET gtcc_paid_%I = %L,
				created_at    = now(),
				updated_at    = now()
			 WHERE id = %L'',
			LOWER(i.moving_expense_type::text),
			i.sum_amt,
			v_ppm_closeout_id
      	);
		EXECUTE v_sql;
	END LOOP;

	--get member paid expenses
	FOR i IN (
		SELECT moving_expense_type, SUM(amount) sum_amt
		  	FROM moving_expenses
		WHERE ppm_shipment_id = p_ppm_shipment_id
		   	AND status = ''APPROVED''
		   	AND paid_with_gtcc = false
		   	AND moving_expense_type != ''STORAGE''
		GROUP BY moving_expense_type)
	LOOP
		v_sql := format(
			''UPDATE ppm_closeouts
			 SET member_paid_%I = %L,
				created_at    = now(),
				updated_at    = now()
			 WHERE id = %L'',
			LOWER(i.moving_expense_type::text),
			i.sum_amt,
			v_ppm_closeout_id
      	);
		EXECUTE v_sql;
	END LOOP;

	--set gtcc total paid expenses
	SELECT
		COALESCE(gtcc_paid_contracted_expense,0) +
		COALESCE(gtcc_paid_oil,0) +
		COALESCE(gtcc_paid_other,0) +
		COALESCE(gtcc_paid_packing_materials,0) +
		COALESCE(gtcc_paid_rental_equipment,0) +
		COALESCE(gtcc_paid_tolls,0) +
		COALESCE(gtcc_paid_weighing_fee,0)
	INTO v_total_gtcc_paid_expenses
	FROM ppm_closeouts
 	WHERE id = v_ppm_closeout_id;

	--set member total paid expenses
	SELECT
		COALESCE(member_paid_contracted_expense,0) +
		COALESCE(member_paid_oil,0) +
		COALESCE(member_paid_other,0) +
		COALESCE(member_paid_packing_materials,0) +
		COALESCE(member_paid_rental_equipment,0) +
		COALESCE(member_paid_tolls,0) +
		COALESCE(member_paid_weighing_fee,0)
	INTO v_total_member_paid_expenses
	FROM ppm_closeouts
 	WHERE id = v_ppm_closeout_id;

	--get sit paid
	SELECT SUM(amount)
	INTO v_gtcc_paid_sit
	FROM moving_expenses
	WHERE ppm_shipment_id = p_ppm_shipment_id
		AND status = ''APPROVED''
	   	AND paid_with_gtcc = true
	   	AND moving_expense_type = ''STORAGE'';

	SELECT SUM(amount)
	INTO v_member_paid_sit
	FROM moving_expenses
	WHERE ppm_shipment_id = p_ppm_shipment_id
		AND status = ''APPROVED''
		AND paid_with_gtcc = false
		AND moving_expense_type = ''STORAGE'';

	--get incentives and advance
	SELECT
		COALESCE(estimated_incentive,0),
		COALESCE(advance_amount_received,0),
		COALESCE(final_incentive,0),
		COALESCE(max_incentive,0),
		ppm_type
	INTO
		v_estimated_incentive,
		v_advance_amount_received,
		v_final_incentive,
		v_max_incentive,
		v_ppm_type
	FROM ppm_shipments
	WHERE id = p_ppm_shipment_id;

	v_gtcc_paid_sit := COALESCE(v_gtcc_paid_sit,0);
	v_member_paid_sit := COALESCE(v_member_paid_sit,0);
	v_max_advance := v_estimated_incentive * .6;
	v_remaining_incentive := v_final_incentive - v_advance_amount_received;
	v_gtcc_disbursement	:= v_total_gtcc_paid_expenses;
	v_member_disbursement := (v_remaining_incentive + v_member_paid_sit) - v_gtcc_disbursement;

	RAISE NOTICE ''Total gtcc paid expenses: %'', v_total_gtcc_paid_expenses;
	RAISE NOTICE ''Total member paid expenses: %'', v_total_member_paid_expenses;
	RAISE NOTICE ''Max advance % = Est incentive % * .6'', v_max_advance, v_estimated_incentive;
	RAISE NOTICE ''Remaining incentive % = Final incentive % - Adv received %'', v_remaining_incentive, v_final_incentive, v_advance_amount_received;
	RAISE NOTICE ''gtcc paid sit: %'', v_gtcc_paid_sit;
	RAISE NOTICE ''member paid sit: %'', v_member_paid_sit;
	RAISE NOTICE ''gtcc disbursement: %'', v_gtcc_disbursement;
	RAISE NOTICE ''member disbursement % = remaining incentive % + member paid sit % - gtcc disbursement %'',v_member_disbursement, v_remaining_incentive, v_member_paid_sit, v_gtcc_disbursement;

	--set ppm_closeout data
	UPDATE ppm_closeouts
	SET max_advance = v_max_advance,
		remaining_incentive = v_remaining_incentive,
		total_gtcc_paid_expenses = v_total_gtcc_paid_expenses,
		total_member_paid_expenses = v_total_member_paid_expenses,
		gtcc_paid_sit = v_gtcc_paid_sit,
		member_paid_sit = v_member_paid_sit,
		member_disbursement = v_member_disbursement,
		gtcc_disbursement = v_gtcc_disbursement
 	WHERE id = v_ppm_closeout_id;

	-- Zero‐fill any unpaid columns for this closeout for consistency with the historic data
	FOR i IN (
		SELECT column_name
		FROM information_schema.columns
	   	WHERE table_schema = ''public''
		 	AND table_name   = ''ppm_closeouts''
		 	AND (column_name LIKE ''gtcc_paid_%''
			OR column_name LIKE ''member_paid_%''))
	LOOP
	  EXECUTE format(
		''UPDATE ppm_closeouts SET
        %I = COALESCE(%I, 0),
				updated_at = now()
		 WHERE id = %L'',
		i.column_name,
		i.column_name,
		v_ppm_closeout_id
	  );
END LOOP;

END;
';
