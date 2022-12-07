import React from 'react';
import { DateField, NumberField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';
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

const WebhookSubscriptionShow = (props) => {
  return (
    /* eslint-disable-next-line react/jsx-props-no-spreading */
    <Show {...props} title={<WebhookSubscriptionShowTitle />}>
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
