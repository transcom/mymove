import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput, required, Toolbar, SaveButton } from 'react-admin';
import { connect } from 'react-redux';

import { selectAdminUser } from 'store/entities/selectors';

const AdminUserEditToolbar = (props) => (
  <Toolbar {...props}>
    <SaveButton />
  </Toolbar>
);

const AdminUserSuperAttribute = () => {
  return (
    <SelectInput
      source="super"
      choices={[
        { id: true, name: 'Yes' },
        { id: false, name: 'No' },
      ]}
      sx={{ width: 256 }}
      disabled
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
      {adminUser?.super && <AdminUserSuperAttribute />}
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
