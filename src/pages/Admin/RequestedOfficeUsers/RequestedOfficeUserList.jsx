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
  useRecordContext,
} from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const RequestedOfficeUserListFilter = (props) => (
  <Filter {...props}>
    <TextInput label="Search by Name/Email" source="search" alwaysOn />
    <TextInput label="Transportation Office" source="offices" alwaysOn />
    <TextInput label="Roles" source="rolesSearch" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'createdAt', order: 'DESC' };

const RolesTextField = (user) => {
  const { roles } = user;

  let roleStr = '';
  for (let i = 0; i < roles.length; i += 1) {
    roleStr += roles[i].roleName;

    if (i < roles.length - 1) {
      roleStr += ', ';
    }
  }

  return roleStr;
};

const RolesField = () => {
  const record = useRecordContext();
  return <div>{RolesTextField(record)}</div>;
};

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
      <ArrayField source="roles" sortable={false} clickable={false} sort={{ field: 'roleName', order: 'DESC' }}>
        <RolesField />
      </ArrayField>
    </Datagrid>
  </List>
);

export default RequestedOfficeUserList;
