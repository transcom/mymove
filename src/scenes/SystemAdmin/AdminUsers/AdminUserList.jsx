import React from 'react';
import { List, Datagrid, TextField, BooleanField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'last_name', order: 'ASC' };

const AdminUserList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25} sort={defaultSort} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="first_name" />
      <TextField source="last_name" />
      <BooleanField source="active" />
    </Datagrid>
  </List>
);

export default AdminUserList;
