import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput } from 'react-admin';

import SaveToolbar from '../Shared/SaveToolbar';

const ClientCertEdit = (props) => (
  <Edit {...props}>
    <SimpleForm toolbar={<SaveToolbar showDeleteBtn />}>
      <TextInput source="id" disabled />
      <TextInput source="userId" disabled label="User Id" />
      <TextInput source="subject" fullWidth />
      <TextInput source="sha256Digest" fullWidth />
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
      <SelectInput
        source="allowPPTAS"
        label="Allow PPTAS"
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
