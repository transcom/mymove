import React from 'react';
import { List, Datagrid, TextField, BooleanField, Filter, TextInput, ReferenceField } from 'react-admin';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const OfficeUserListFilter = props => (
  <Filter {...props}>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'last_name', order: 'ASC' };

const OfficeUserList = props => (
  <List
    {...props}
    pagination={<AdminPagination />}
    perPage={25}
    bulkActionButtons={false}
    sort={defaultSort}
    filters={<OfficeUserListFilter />}
  >
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="first_name" />
      <TextField source="last_name" />
      <ReferenceField label="Transportation Office" source="transportation_office_id" reference="offices">
        <TextField source="name" />
      </ReferenceField>
      <BooleanField source="active" />
    </Datagrid>
  </List>
);

export default OfficeUserList;
