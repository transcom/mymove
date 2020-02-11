UPDATE reimbursements r
SET method_of_receipt = 'OTHER_DD'
WHERE
	method_of_receipt = 'GTCC' AND
	r.id in (
		SELECT advance_id from personally_procured_moves
	);
