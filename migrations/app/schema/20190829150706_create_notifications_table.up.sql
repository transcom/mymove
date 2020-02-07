CREATE TABLE notifications
(
    id                uuid PRIMARY KEY,
    service_member_id uuid
        CONSTRAINT service_member_id___fk
            REFERENCES service_members,
    ses_message_id    text                     NOT NULL,
    notification_type text                     NOT NULL,
    created_at        timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);
