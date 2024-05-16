import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { roleTypes } from 'constants/userRoles';

const OfficeUserCreate = () => {
  const validateForm = (values) => {
    const errors = {};
    if (!values.firstName) {
      errors.firstName = 'First name is required';
    }
    if (!values.lastName) {
      errors.lastName = 'Last name is required';
    }
    if (!values.email) {
      errors.email = 'Email is required';
    }

    if (!values.telephone) {
      errors.telephone = 'Telephone is required.';
    } else if (!values.telephone.match(/^[2-9]\d{2}-\d{3}-\d{4}$/)) {
      errors.telephone = 'Invalid phone number, should be 000-000-0000.';
    }

    if (!values.roles?.length) {
      errors.roles = 'Role(s) are required.';
    } else if (
      values.roles.find((role) => role.roleType === roleTypes.TIO) &&
      values.roles.find((role) => role.roleType === roleTypes.TOO)
    ) {
      errors.roles = 'A user must not have both TOO and TIO role.';
    }

    if (!values.transportationOfficeId) {
      errors.transportationOfficeId = 'Transportation Office is required';
    }

    return errors;
  };

  return (
    <Create>
      <SimpleForm
        sx={{ '& .MuiInputBase-input': { width: 232 } }}
        reValidateMode="onBlur"
        mode="onBlur"
        validate={validateForm}
      >
        <TextInput source="firstName" validate={required()} />
        <TextInput source="middleInitials" />
        <TextInput source="lastName" validate={required()} />
        <TextInput source="email" validate={required()} />
        <TextInput source="telephone" validate={required()} />
        <RolesPrivilegesCheckboxInput source="roles" validate={required()} />
        <ReferenceInput
          label="Transportation Office"
          reference="offices"
          source="transportationOfficeId"
          perPage={500}
          validate={required()}
        >
          <AutocompleteInput optionText="name" validate={required()} sx={{ width: 256 }} />
        </ReferenceInput>
      </SimpleForm>
    </Create>
  );
};

export default OfficeUserCreate;
