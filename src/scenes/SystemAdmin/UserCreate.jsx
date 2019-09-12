import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required, regex } from 'react-admin';

const phoneValidators = [required(), regex(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Invalid phone number, should be 000-000-0000')];

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
