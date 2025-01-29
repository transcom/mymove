import React from 'react';
import {
  ArrayField,
  Datagrid,
  DateField,
  Filter,
  List,
  ReferenceField,
  TextField,
  TextInput,
  TopToolbar,
} from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar
const ListActions = () => {
  return <TopToolbar />;
};

const RejectedOfficeUserListFilter = () => (
  <Filter>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'createdAt', order: 'DESC' };

const RejectedOfficeUserList = () => (
  <List
    pagination={<AdminPagination />}
    perPage={25}
    sort={defaultSort}
    filters={<RejectedOfficeUserListFilter />}
    actions={<ListActions />}
  >
    <Datagrid bulkActionButtons={false} rowClick="show" data-testid="rejected-office-user-fields">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="firstName" />
      <TextField source="lastName" />
      <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" link={false}>
        <TextField source="name" />
      </ReferenceField>
      <TextField source="status" />
      <TextField source="rejectionReason" label="Reason for rejection" />
      <DateField showTime source="rejectedOn" label="Rejected date" />
      <ArrayField source="roles" label="Requested Roles">
        <Datagrid bulkActionButtons={false} headerHeight="0" sx={{ paddingTop: 0, paddingBottom: 0 }}>
          <TextField source="roleName" label="" />
        </Datagrid>
      </ArrayField>
    </Datagrid>
  </List>
);

export default RejectedOfficeUserList;
