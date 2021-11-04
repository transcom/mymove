import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput, Toolbar, DeleteButton, SaveButton } from 'react-admin';

const ClientCertEditToolbar = (props) => (
  <Toolbar {...props}>
    <DeleteButton />
    <SaveButton />
  </Toolbar>
);

const ClientCertEdit = (props) => (
  <Edit {...props}>
    <SimpleForm toolbar={<ClientCertEditToolbar />}>
      <TextInput source="id" disabled />
      <TextInput source="subject" disabled />
      <TextInput source="sha256Digest" disabled />
      <SelectInput
        source="allowDpsAuthAPI"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
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
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
    </SimpleForm>
  </Edit>
);

export default ClientCertEdit;
