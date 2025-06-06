import React from 'react';
import { Create, SimpleForm, TextInput, SelectInput, required } from 'react-admin';

import SaveToolbar from '../Shared/SaveToolbar';

import { WEBHOOK_SUBSCRIPTION_STATUS } from 'shared/constants';

const WebhookSubscriptionCreate = () => (
  <Create>
    <SimpleForm
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
      toolbar={<SaveToolbar />}
    >
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
        sx={{ width: 256 }}
      />
    </SimpleForm>
  </Create>
);

export default WebhookSubscriptionCreate;
