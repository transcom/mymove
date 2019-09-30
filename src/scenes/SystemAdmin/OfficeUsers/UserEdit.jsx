import React from 'react';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import { Edit, SimpleForm, TextInput, DisabledInput, required, Toolbar, SaveButton } from 'react-admin';

const UserEditToolbar = props => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const UserEdit = props => (
  <Edit {...props}>
    <SimpleForm toolbar={<UserEditToolbar />}>
      <DisabledInput source="id" />
      <DisabledInput source="email" />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="middle_initials" />
      <TextInput source="last_name" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <DisabledInput source="disabled" label="Deactivated" />
      <DisabledInput source="created_at" />
      <DisabledInput source="updated_at" />
    </SimpleForm>
  </Edit>
);

export default UserEdit;
