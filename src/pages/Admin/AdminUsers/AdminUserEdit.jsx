import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput, required, Toolbar, SaveButton, useRecordContext } from 'react-admin';
import { connect } from 'react-redux';

import { selectAdminUser } from 'store/entities/selectors';

const AdminUserEditToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const AdminUserSuperAttribute = ({ adminUser }) => {
  const record = useRecordContext();
  // Hide the input so the super admin can't un-super themselves
  if (record.id === adminUser.id) {
    return null;
  }
  return (
    <SelectInput
      source="super"
      choices={[
        { id: true, name: 'Yes' },
        { id: false, name: 'No' },
      ]}
      sx={{ width: 256 }}
    />
  );
};

const AdminUserEdit = ({ adminUser }) => (
  <Edit>
    <SimpleForm
      toolbar={<AdminUserEditToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
    >
      <TextInput source="id" disabled />
      <TextInput source="userId" label="User Id" disabled />
      <TextInput source="email" disabled />
      <TextInput source="firstName" validate={required()} />
      <TextInput source="lastName" validate={required()} />
      <SelectInput
        source="active"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      {adminUser?.super && <AdminUserSuperAttribute adminUser={adminUser} />}
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
    </SimpleForm>
  </Edit>
);

function mapStateToProps(state) {
  return {
    adminUser: selectAdminUser(state),
  };
}

export default connect(mapStateToProps)(AdminUserEdit);
