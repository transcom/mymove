import React from 'react';
import {
  Datagrid,
  DateField,
  Filter,
  List,
  ReferenceField,
  TextField,
  TextInput,
  TopToolbar,
  useRecordContext,
} from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const RejectedOfficeUserShowRoles = () => {
  const officeUser = useRecordContext();
  if (!officeUser?.roles) return <p>This user has not rejected any roles.</p>;

  const uniqueRoleNamesList = [];
  const rejectedRoles = officeUser.roles;
  for (let i = 0; i < rejectedRoles.length; i += 1) {
    if (!uniqueRoleNamesList.includes(rejectedRoles[i].roleName)) {
      uniqueRoleNamesList.push(rejectedRoles[i].roleName);
    }
  }

  return <p>{uniqueRoleNamesList.join(', ')}</p>;
};

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
      <RejectedOfficeUserShowRoles sortable={false} source="roles" label="Rejected Roles" />
    </Datagrid>
  </List>
);

export default RejectedOfficeUserList;
