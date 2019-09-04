import React from 'react';
import { List, Datagrid, TextField } from 'react-admin';
import AdminPagination from './AdminPagination';

const AccessCodeList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={500}>
    <Datagrid>
      <TextField source="id" />
      <TextField source="code" />
      <TextField source="move_type" />
      <TextField source="locator" />
    </Datagrid>
  </List>
);

export default AccessCodeList;
