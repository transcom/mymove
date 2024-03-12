import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';

const OfficeUserCreate = () => {
  return (
    <Create>
      <SimpleForm sx={{ '& .MuiInputBase-input': { width: 232 } }} mode="onBlur" reValidateMode="onBlur">
        <TextInput source="firstName" validate={required()} />
        <TextInput source="middleInitials" />
        <TextInput source="lastName" validate={required()} />
        <TextInput source="email" validate={required()} />
        <TextInput source="telephone" validate={phoneValidators} />
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
