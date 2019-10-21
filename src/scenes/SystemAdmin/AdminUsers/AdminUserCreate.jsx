import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, SelectInput, required } from 'react-admin';

const AdminUserCreate = props => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="email" validate={required()} />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="last_name" validate={required()} />
      <ReferenceInput label="Organization" reference="organizations" source="organization_id" perPage={500}>
        <AutocompleteInput optionText="name" />
      </ReferenceInput>
      <SelectInput
        source="role"
        choices={[{ id: 'SYSTEM_ADMIN', name: 'System Admin' }, { id: 'PROGRAM_ADMIN', name: 'Program Admin' }]}
      />
    </SimpleForm>
  </Create>
);

export default AdminUserCreate;
