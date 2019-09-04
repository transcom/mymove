import React from 'react';
import { List, EmailField, Datagrid, TextField } from 'react-admin';
import AdminPagination from './AdminPagination';

const UserList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25}>
    <Datagrid rowClick="show">
      <EmailField source="email" />
      <TextField source="first_name" />
      <TextField source="id" />
      <TextField source="last_name" />
    </Datagrid>
  </List>
);

export default UserList;
