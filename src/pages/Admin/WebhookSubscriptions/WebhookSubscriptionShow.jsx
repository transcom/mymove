import React from 'react';
import { DateField, NumberField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';

const WebhookSubscriptionShowTitle = () => {
  const record = useRecordContext();
  return <span>{`Webhook Subscription ID: ${record.id}`}</span>;
};

const WebhookSubscriptionShow = () => {
  return (
    <Show title={<WebhookSubscriptionShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField label="Subscriber Id" source="subscriberId" />
        <TextField source="eventKey" />
        <TextField source="callbackUrl" />
        <NumberField source="severity" />
        <TextField source="status" />
        <DateField source="updatedAt" showTime />
        <DateField source="createdAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default WebhookSubscriptionShow;
