import React from 'react';
import { List, Datagrid, TextField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'name', order: 'ASC' };

const OfficeList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25} sort={defaultSort}>
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
