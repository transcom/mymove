import React from 'react';
import { Create, SimpleForm, TextInput, ReferenceInput, AutocompleteInput, required } from 'react-admin';

import SaveToolbar from '../Shared/SaveToolbar';

const AdminUserCreate = () => (
  <Create>
    <SimpleForm
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
      toolbar={<SaveToolbar />}
    >
      <TextInput source="email" validate={required()} />
      <TextInput source="firstName" validate={required()} />
      <TextInput source="lastName" validate={required()} />
      <ReferenceInput label="Organization" reference="organizations" source="organizationId" perPage={500}>
        <AutocompleteInput optionText="name" validate={required()} sx={{ width: 256 }} />
      </ReferenceInput>
    </SimpleForm>
  </Create>
);

export default AdminUserCreate;
