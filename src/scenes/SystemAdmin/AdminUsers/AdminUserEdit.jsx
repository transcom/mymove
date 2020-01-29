import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput, required, Toolbar, SaveButton } from 'react-admin';

const AdminUserEditToolbar = props => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const AdminUserEdit = props => (
  <Edit {...props}>
    <SimpleForm toolbar={<AdminUserEditToolbar />}>
      <TextInput source="id" disabled />
      <TextInput source="email" disabled />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="last_name" validate={required()} />
      <SelectInput
        source="active"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
      />
      <TextInput source="created_at" disabled />
      <TextInput source="updated_at" disabled />
    </SimpleForm>
  </Edit>
);

export default AdminUserEdit;
