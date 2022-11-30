import React from 'react';
import { Create, SimpleForm, TextInput, SelectInput, required } from 'react-admin';

import { WEBHOOK_SUBSCRIPTION_STATUS } from 'shared/constants';

const WebhookSubscriptionCreate = () => (
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  <Create>
    <SimpleForm>
      <TextInput label="Subscriber Id" source="subscriberId" validate={required()} />
      <TextInput source="eventKey" validate={required()} />
      <TextInput source="callbackUrl" validate={required()} />
      <SelectInput
        source="status"
        choices={[
          { id: WEBHOOK_SUBSCRIPTION_STATUS.ACTIVE, name: 'Active' },
          { id: WEBHOOK_SUBSCRIPTION_STATUS.DISABLED, name: 'Disabled' },
          { id: WEBHOOK_SUBSCRIPTION_STATUS.FAILING, name: 'Failing' },
        ]}
      />
    </SimpleForm>
  </Create>
);

export default WebhookSubscriptionCreate;
