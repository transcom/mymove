import React from 'react';
import { Show, SimpleShowLayout, TextField, DateField, NumberField } from 'react-admin';
import PropTypes from 'prop-types';

const WebhookSubscriptionShowTitle = ({ record }) => {
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
        <TextField source="subscriberId" />
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
