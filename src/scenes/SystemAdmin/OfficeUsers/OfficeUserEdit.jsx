import React from 'react';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import { Edit, SimpleForm, TextInput, SelectInput, required, Toolbar, SaveButton } from 'react-admin';

const OfficeUserEditToolbar = props => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const OfficeUserEdit = props => (
  <Edit {...props}>
    <SimpleForm toolbar={<OfficeUserEditToolbar />}>
      <TextInput source="id" disabled />
      <TextInput source="email" disabled />
      <TextInput source="first_name" validate={required()} />
      <TextInput source="middle_initials" />
      <TextInput source="last_name" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
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

export default OfficeUserEdit;
