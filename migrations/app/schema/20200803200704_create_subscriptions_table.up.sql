CREATE TYPE webhook_subscriptions_status AS ENUM (
    'ACTIVE',
    'DISABLED'
);

CREATE TABLE webhook_subscriptions
(
    id                  uuid                          PRIMARY KEY            NOT NULL,
    subscriber_id       uuid                          REFERENCES contractors NOT NULL,
    status              webhook_subscriptions_status                         DEFAULT 'ACTIVE',
    event_key           text                                                 NOT NULL,
    callback_url        text                                                 NOT NULL,
    created_at          timestamp WITH TIME ZONE                             NOT NULL,
    updated_at          timestamp WITH TIME ZONE                             NOT NULL
);

