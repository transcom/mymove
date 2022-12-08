import React from 'react';
import { DateField, NumberField, Show, SimpleShowLayout, TextField } from 'react-admin';
import PropTypes from 'prop-types';

const WebhookSubscriptionShowTitle = () => {
  const record = useRecordContext();
  return <span>{`Webhook Subscription ID: ${record.id}`}</span>;
};

WebhookSubscriptionShowTitle.propTypes = {
  record: PropTypes.shape({
    id: PropTypes.string,
  }),
};

WebhookSubscriptionShowTitle.defaultProps = {
  record: {
    id: '',
  },
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
        <DateField source="updatedAt" showTime addLabel />
        <DateField source="createdAt" showTime addLabel />
      </SimpleShowLayout>
    </Show>
  );
};

export default WebhookSubscriptionShow;
