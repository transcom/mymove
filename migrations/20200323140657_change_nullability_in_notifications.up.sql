-- Change nullability
Alter table notifications
    ALTER COLUMN service_member_id SET NOT NULL;
