import React from 'react';
import { List, Datagrid, TextField, Filter, TextInput } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'name', order: 'ASC' };

const OfficeFilter = (props) => (
  <Filter {...props}>
    <TextInput label="Search by Office Name" source="q" resettable alwaysOn />
  </Filter>
);

const OfficeList = () => (
  <List filters={<OfficeFilter />} pagination={<AdminPagination />} perPage={25} sort={defaultSort}>
    <Datagrid bulkActionButtons={false}>
      <TextField source="id" />
      <TextField source="name" />
      <TextField source="latitude" />
      <TextField source="longitude" />
      <TextField source="gbloc" />
    </Datagrid>
  </List>
);

export default OfficeList;
