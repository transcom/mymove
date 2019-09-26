import React from 'react';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

const UserCreate = props => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="first_name" validate={required()} />
      <TextInput source="middle_initial" />
      <TextInput source="last_name" validate={required()} />
      <TextInput source="email" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <ReferenceInput label="Transportation Office" reference="offices" source="transportation_office_id" perPage={500}>
        <AutocompleteInput optionText="name" />
      </ReferenceInput>
    </SimpleForm>
  </Create>
);

export default UserCreate;
