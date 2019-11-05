import React from 'react';
import { Edit, SimpleForm, TextInput, DisabledInput, SelectInput, required, Toolbar, SaveButton } from 'react-admin';

const AdminUserEditToolbar = props => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const AdminUserEdit = props => (
  <Edit {...props}>
    <SimpleForm toolbar={<AdminUserEditToolbar />}>
      <DisabledInput source="id" />
      <DisabledInput source="email" />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="last_name" validate={required()} />
      <SelectInput source="active" choices={[{ id: true, name: 'Yes' }, { id: false, name: 'No' }]} />
      <DisabledInput source="created_at" />
      <DisabledInput source="updated_at" />
    </SimpleForm>
  </Edit>
);

export default AdminUserEdit;
