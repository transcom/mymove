import React from 'react';
import { List, Datagrid, TextField, BooleanField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'email', order: 'ASC' };

const CustomerUserList = (props) => (
  <List {...props} pagination={<AdminPagination />} perPage={25} sort={defaultSort} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <BooleanField source="active" />
      <TextField source="createAt" />
    </Datagrid>
  </List>
);

export default CustomerUserList;
