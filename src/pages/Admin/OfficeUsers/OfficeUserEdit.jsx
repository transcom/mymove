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
import { roleTypes } from 'constants/userRoles';

const OfficeUserEditToolbar = (props) => {
  return (
    <Toolbar {...props}>
      <SaveButton />
    </Toolbar>
  );
};
const validateForm = (values) => {
  const errors = {};
  if (!values.firstName) {
    errors.firstName = 'You must enter a first name.';
  }
  if (!values.lastName) {
    errors.lastName = 'You must enter a last name.';
  }
  if (!values.email) {
    errors.email = 'You must enter an email.';
  }

  if (!values.telephone) {
    errors.telephone = 'You must enter a telephone number.';
  } else if (!values.telephone.match(/^[2-9]\d{2}-\d{3}-\d{4}$/)) {
    errors.telephone = 'Invalid phone number, should be 000-000-0000.';
  }

  if (!values.roles?.length) {
    errors.roles = 'You must select at least one role.';
  } else if (
    values.roles.find((role) => role.roleType === roleTypes.TIO) &&
    values.roles.find((role) => role.roleType === roleTypes.TOO)
  ) {
    errors.roles =
      'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.';
  }

  if (!values.transportationOfficeId) {
    errors.transportationOfficeId = 'You must select a transportation office.';
  }

  return errors;
};

const OfficeUserEdit = () => (
  <Edit>
    <SimpleForm
      toolbar={<OfficeUserEditToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
      validate={validateForm}
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
