import React from 'react';
import { List, Datagrid, TextField } from 'react-admin';
import AdminPagination from './AdminPagination';

const OfficeList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={500}>
    <Datagrid>
      <TextField source="id" />
      <TextField source="name" />
      <TextField source="latitude" />
      <TextField source="longitude" />
      <TextField source="gbloc" />
    </Datagrid>
  </List>
);

export default OfficeList;
