import React, { useEffect, useState } from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

import { RolesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesCheckboxes';
import { PrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/ElevatedPrivilegeCheckboxes';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';

const OfficeUserCreate = () => {
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
    <Create>
      <SimpleForm sx={{ '& .MuiInputBase-input': { width: 232 } }} mode="onBlur" reValidateMode="onBlur">
        <TextInput source="firstName" validate={required()} />
        <TextInput source="middleInitials" />
        <TextInput source="lastName" validate={required()} />
        <TextInput source="email" validate={required()} />
        <TextInput source="telephone" validate={phoneValidators} />
        <RolesCheckboxInput source="roles" validate={required()} onChange={validatePrivileges} />
        <PrivilegesCheckboxInput source="elevatedPrivileges" disabled={isDisabledPrivileges} />
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
