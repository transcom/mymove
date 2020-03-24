-- Change nullability
ALTER TABLE notifications
    ALTER COLUMN service_member_id SET NOT NULL;
