import React from 'react';
import { BooleanField, Datagrid, List, TextField } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'last_name', order: 'ASC' };

const AdminUserList = () => (
  <List pagination={<AdminPagination />} perPage={25} sort={defaultSort} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="firstName" />
      <TextField source="lastName" />
      <TextField source="userId" label="User Id" />
      <BooleanField source="active" />
    </Datagrid>
  </List>
);

export default AdminUserList;
