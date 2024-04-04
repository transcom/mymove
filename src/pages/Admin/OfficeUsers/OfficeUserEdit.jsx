import React from 'react';
import {
  Edit,
  SimpleForm,
  TextInput,
  SelectInput,
  required,
  Toolbar,
  SaveButton,
  AutocompleteInput,
  ReferenceInput,
} from 'react-admin';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';

const OfficeUserEditToolbar = (props) => {
  return (
    <Toolbar {...props}>
      <SaveButton />
    </Toolbar>
  );
};

const OfficeUserEdit = () => (
  <Edit>
    <SimpleForm
      toolbar={<OfficeUserEditToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
    >
      <TextInput source="id" disabled />
      <TextInput source="userId" label="User Id" disabled />
      <TextInput source="email" disabled />
      <TextInput source="firstName" validate={required()} />
      <TextInput source="middleInitials" />
      <TextInput source="lastName" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <SelectInput
        source="active"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <RolesPrivilegesCheckboxInput source="roles" validate={required()} />
      <ReferenceInput
        label="Transportation Office"
        reference="offices"
        source="transportationOfficeId"
        perPage={500}
        validate={required()}
      >
        <AutocompleteInput optionText="name" sx={{ width: 256 }} />
      </ReferenceInput>
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
    </SimpleForm>
  </Edit>
);

export default OfficeUserEdit;
