import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput, required, Toolbar, SaveButton } from 'react-admin';

const AdminUserEditToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const AdminUserEdit = () => (
  <Edit>
    <SimpleForm toolbar={<AdminUserEditToolbar />}>
      <TextInput source="id" disabled />
      <TextInput source="userId" label="User Id" disabled />
      <TextInput source="email" disabled />
      <TextInput source="firstName" validate={required()} />
      <TextInput source="lastName" validate={required()} />
      <SelectInput
        source="active"
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

export default AdminUserEdit;
