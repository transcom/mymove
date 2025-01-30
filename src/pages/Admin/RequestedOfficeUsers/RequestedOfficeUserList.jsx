import React from 'react';
import { Datagrid, DateField, Filter, List, ReferenceField, TextField, TextInput } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const RequestedOfficeUserListFilter = (props) => (
  <Filter {...props}>
    <TextInput label="Search by Name/Email" source="search" alwaysOn />
    <TextInput label="Transportation Office" source="offices" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'createdAt', order: 'DESC' };

const RequestedOfficeUserList = () => (
  <List pagination={<AdminPagination />} perPage={25} sort={defaultSort} filters={<RequestedOfficeUserListFilter />}>
    <Datagrid bulkActionButtons={false} rowClick="show" data-testid="requested-office-user-fields">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="firstName" />
      <TextField source="lastName" />
      <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" link={false}>
        <TextField source="name" />
      </ReferenceField>
      <TextField source="status" />
      <DateField showTime source="createdAt" label="Requested on" />
    </Datagrid>
  </List>
);

export default RequestedOfficeUserList;
