UPDATE moves m
SET m.cancel_reason = 'Terminated for Cause (TCN-0001)',
    mto.status  = 'CANCELED'
FROM mto_shipments mto
WHERE m.id = mto.move_id
and moves.locator in (

);