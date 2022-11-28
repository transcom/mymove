import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

const AdminUserCreate = () => (
  <Create>
    <SimpleForm>
      <TextInput source="email" validate={required()} />
      <TextInput source="firstName" validate={required()} />
      <TextInput source="lastName" validate={required()} />
      <ReferenceInput
        label="Organization"
        reference="organizations"
        source="organizationId"
        perPage={500}
        validate={required()}
      >
        <AutocompleteInput optionText="name" />
      </ReferenceInput>
    </SimpleForm>
  </Create>
);

export default AdminUserCreate;
