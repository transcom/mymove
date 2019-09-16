import React from 'react';
import { Edit, SimpleForm, TextInput, DisabledInput } from 'react-admin';

const UserEdit = props => (
  <Edit {...props}>
    <SimpleForm>
      <DisabledInput source="id" />
      <DisabledInput source="email" />
      <TextInput source="first_name" />
      <TextInput source="middle_initials" />
      <TextInput source="last_name" />
      <TextInput source="telephone" />
    </SimpleForm>
  </Edit>
);

export default UserEdit;
