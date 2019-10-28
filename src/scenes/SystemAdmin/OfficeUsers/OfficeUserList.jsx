import React from 'react';
import { List, Datagrid, TextField, BooleanField, Filter, TextInput } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const OfficeUserListFilter = props => (
  <Filter {...props}>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const OfficeUserList = props => (
  <List
    {...props}
    pagination={<AdminPagination />}
    perPage={25}
    bulkActionButtons={false}
    filters={<OfficeUserListFilter />}
  >
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
