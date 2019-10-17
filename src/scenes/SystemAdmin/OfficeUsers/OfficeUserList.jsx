import React from 'react';
import { List, Datagrid, TextField, BooleanField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const OfficeUserList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="first_name" />
      <TextField source="last_name" />
      <BooleanField source="deactivated" />
    </Datagrid>
  </List>
);

export default OfficeUserList;
