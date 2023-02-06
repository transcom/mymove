import React from 'react';
import { Create, SimpleForm, TextInput, SelectInput, required } from 'react-admin';

const ClientCertCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="subject" validate={required()} />
      <TextInput source="sha256Digest" validate={required()} />
      <TextInput source="user_id" validate={required()} />
      <SelectInput
        source="allowOrdersAPI"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowAirForceOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowAirForceOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowArmyOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowArmyOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowCoastGuardOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowCoastGuardOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowMarineCorpsOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowMarineCorpsOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowNavyOrdersRead"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowNavyOrdersWrite"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <SelectInput
        source="allowPrime"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
    </SimpleForm>
  </Create>
);

export default ClientCertCreate;
