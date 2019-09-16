import React from 'react';
import { phoneValidators } from './form_validators';
import { Edit, SimpleForm, TextInput, DisabledInput, required } from 'react-admin';

const UserEdit = props => (
  <Edit {...props}>
    <SimpleForm>
      <DisabledInput source="id" />
      <DisabledInput source="email" />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="middle_initials" />
      <TextInput source="last_name" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <DisabledInput source="disabled" />
      <DisabledInput source="created_at" />
      <DisabledInput source="updated_at" />
    </SimpleForm>
  </Edit>
);

export default UserEdit;
