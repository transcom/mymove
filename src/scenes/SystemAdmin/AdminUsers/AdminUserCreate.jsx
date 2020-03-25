import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

const AdminUserCreate = props => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="email" validate={required()} />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="last_name" validate={required()} />
      <ReferenceInput label="Organization" reference="organizations" source="organization_id" perPage={500}>
        <AutocompleteInput optionText="name" />
      </ReferenceInput>
    </SimpleForm>
  </Create>
);

export default AdminUserCreate;
