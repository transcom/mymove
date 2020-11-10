import React from 'react';
import { RolesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesCheckboxes';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

const OfficeUserCreate = (props) => (
  <Create {...props}>
    <SimpleForm>
      <TextInput source="firstName" validate={required()} />
      <TextInput source="middle_initial" />
      <TextInput source="lastName" validate={required()} />
      <TextInput source="email" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <RolesCheckboxInput source="roles" validate={required()} />
      <ReferenceInput label="Transportation Office" reference="offices" source="transportationOfficeId" perPage={500}>
        <AutocompleteInput optionText="name" />
      </ReferenceInput>
    </SimpleForm>
  </Create>
);

export default OfficeUserCreate;
