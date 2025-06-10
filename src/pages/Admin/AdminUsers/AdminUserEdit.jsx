import React from 'react';
import { Edit, SimpleForm, TextInput, SelectInput, required } from 'react-admin';
import { connect } from 'react-redux';

import SaveToolbar from '../Shared/SaveToolbar';

import { selectAdminUser } from 'store/entities/selectors';

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

const validateAdminuser = (values) => {
  const errors = {};
  if (!values.email) {
    errors.email = 'Email is required';
  } else if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(values.email)) {
    errors.email = 'Must be a valid email';
  }
  if (!values.firstName) {
    errors.firstName = 'First name is required';
  }
  if (!values.lastName) {
    errors.lastName = 'Last name is required';
  }
  return errors;
};

const AdminUserEdit = ({ adminUser }) => (
  <Edit>
    <SimpleForm
      toolbar={<SaveToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
      validate={validateAdminuser}
    >
      <TextInput source="id" disabled />
      <TextInput source="userId" label="User Id" disabled />
      <TextInput source="email" />
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
