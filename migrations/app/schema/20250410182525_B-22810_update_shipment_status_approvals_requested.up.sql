--service item request
update mto_shipments
   set status = 'APPROVALS_REQUESTED',
       updated_at = now()
where id in (
select b.id
from moves a
join mto_shipments b on a.id = b.move_id
join mto_service_items c on b.id = c.mto_shipment_id
where a.status = 'APPROVALS REQUESTED'
and b.status = 'APPROVED'
and c.status = 'SUBMITTED'
and b.shipment_type != 'PPM');

--dest address update
update mto_shipments
   set status = 'APPROVALS_REQUESTED',
       updated_at = now()
where id in (
select b.id
from moves a
join mto_shipments b on a.id = b.move_id
join shipment_address_updates c on b.id = c.shipment_id
where a.status = 'APPROVALS REQUESTED'
and b.status = 'APPROVED'
and c.status = 'REQUESTED'
and b.shipment_type != 'PPM');

--sit extension
update mto_shipments
   set status = 'APPROVALS_REQUESTED',
       updated_at = now()
where id in (
select b.id
from moves a
join mto_shipments b on a.id = b.move_id
join sit_extensions c on b.id = c.mto_shipment_id
where a.status = 'APPROVALS REQUESTED'
and b.status = 'APPROVED'
and c.status = 'PENDING'
and b.shipment_type != 'PPM');