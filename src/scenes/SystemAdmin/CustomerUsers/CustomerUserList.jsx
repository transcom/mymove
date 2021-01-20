import React from 'react';
import { List, Datagrid, TextField, BooleanField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'loginGovEmail', order: 'ASC' };

const CustomerUserList = (props) => (
  <List {...props} pagination={<AdminPagination />} perPage={25} sort={defaultSort} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="loginGovEmail" label="Email" />
      <BooleanField source="active" />
      <TextField source="createdAt" />
    </Datagrid>
  </List>
);

export default CustomerUserList;
