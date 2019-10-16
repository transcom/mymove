import React from 'react';
import { List, Datagrid, TextField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';
import TitleizedField from 'scenes/SystemAdmin/shared/TitleizedField';

const ElectronicOrderList = props => (
  <List {...props} pagination={<AdminPagination />} perPage={25}>
    <Datagrid>
      <TextField source="id" />
      <TitleizedField source="issuer" />
      <TextField source="created_at" />
    </Datagrid>
  </List>
);

export default ElectronicOrderList;
