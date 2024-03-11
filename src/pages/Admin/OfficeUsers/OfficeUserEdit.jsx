import React, { useEffect, useState } from 'react';
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

import { RolesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesCheckboxes';
import { PrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/ElevatedPrivilegeCheckboxes';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';

const OfficeUserEditToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const OfficeUserEdit = () => {
  const [isDisabledPrivileges, setIsDisabledPrivileges] = useState(false);

  const validatePrivileges = (input) => {
    for (let i = 0; i < input?.length; i += 1) {
      if (input[i] === 'customer' || input[i] === 'contracting_officer') {
        setIsDisabledPrivileges(true);
        return;
      }
    }
    setIsDisabledPrivileges(false);
  };

  useEffect(() => {
    validatePrivileges();
  }, []);
  return (
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
        <RolesCheckboxInput source="roles" validate={required()} onChange={validatePrivileges} />
        <PrivilegesCheckboxInput source="privileges" disabled={isDisabledPrivileges} />
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
};

export default OfficeUserEdit;
