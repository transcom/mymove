import React from 'react';
import { Edit, required, SaveButton, SelectInput, SimpleForm, TextInput, Toolbar } from 'react-admin';

import { WEBHOOK_SUBSCRIPTION_STATUS } from 'shared/constants';

const WebhookSubscriptionEditToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const WebhookSubscriptionEdit = (props) => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Edit {...props}>
    <SimpleForm toolbar={<WebhookSubscriptionEditToolbar />}>
      <TextInput source="id" disabled />
      <TextInput label="Subscriber Id" source="subscriberId" validate={required()} />
      <TextInput source="eventKey" validate={required()} />
      <TextInput source="callbackUrl" validate={required()} />
      <SelectInput
        source="severity"
        choices={[
          { id: 0, name: '0' },
          { id: 1, name: '1' },
          { id: 2, name: '2' },
          { id: 3, name: '3' },
          { id: 4, name: '4' },
        ]}
      />
      <SelectInput
        source="status"
        choices={[
          { id: WEBHOOK_SUBSCRIPTION_STATUS.ACTIVE, name: 'Active' },
          { id: WEBHOOK_SUBSCRIPTION_STATUS.DISABLED, name: 'Disabled' },
          { id: WEBHOOK_SUBSCRIPTION_STATUS.FAILING, name: 'Failing' },
        ]}
      />
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
    </SimpleForm>
  </Edit>
);

export default WebhookSubscriptionEdit;
