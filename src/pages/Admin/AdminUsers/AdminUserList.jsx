import React from 'react';
import { BooleanField, Datagrid, List, TextField } from 'react-admin';
import { connect } from 'react-redux';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import { selectAdminUser } from 'store/entities/selectors';

const defaultSort = { field: 'last_name', order: 'ASC' };

const AdminUserList = ({ adminUser }) => (
  <List pagination={<AdminPagination />} perPage={25} sort={defaultSort}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="firstName" />
      <TextField source="lastName" />
      <TextField source="userId" label="User Id" />
      <BooleanField source="active" />
      {adminUser?.super && <BooleanField source="super" label="Super Admin" />}
    </Datagrid>
  </List>
);

function mapStateToProps(state) {
  return {
    adminUser: selectAdminUser(state),
  };
}

export default connect(mapStateToProps)(AdminUserList);
