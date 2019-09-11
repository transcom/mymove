import React from 'react';
import { Create, SimpleForm, TextInput, required, regex } from 'react-admin';

const phoneValidators = [required(), regex(/^[2-9]\d{2}-\d{3}-\d{4}$/, 'Invalid phone number, should be 000-000-0000')];
const uuidValidators = [
  required(),
  regex(
    /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/,
    'Invalid uuid format, should be 12345678-1234-1234-1234-1234567890ab',
  ),
];

const UserCreate = props => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="first_name" validate={required()} />
      <TextInput source="middle_initial" />
      <TextInput source="last_name" validate={required()} />
      <TextInput source="email" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <TextInput source="transportation_office_id" validate={uuidValidators} />
    </SimpleForm>
  </Create>
);

export default UserCreate;
