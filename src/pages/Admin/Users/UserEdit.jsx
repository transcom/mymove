import React from 'react';
import { Edit, SaveButton, SelectInput, SimpleForm, TextInput, Toolbar } from 'react-admin';

const UserEditToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const UserEdit = () => (
  <Edit>
    <SimpleForm
      toolbar={<UserEditToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
    >
      <TextInput source="id" disabled />
      <TextInput source="loginGovEmail" disabled />
      <SelectInput
        source="active"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <SelectInput
        source="revokeAdminSession"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <SelectInput
        source="revokeOfficeSession"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <SelectInput
        source="revokeMilSession"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
    </SimpleForm>
  </Edit>
);

export default UserEdit;
