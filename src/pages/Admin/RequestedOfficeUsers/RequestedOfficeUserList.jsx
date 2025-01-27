import { React } from 'react';
import {
  Datagrid,
  DateField,
  Filter,
  List,
  ReferenceField,
  TextField,
  TopToolbar,
  ArrayField,
  SearchInput,
  useRecordContext,
} from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar
const ListActions = () => {
  return <TopToolbar />;
};

const RequestedOfficeUserListFilter = () => (
  <Filter>
    <SearchInput source="search" alwaysOn />
    <SearchInput source="transportationOfficeSearch" alwaysOn resettable placeholder="Transportation Office" />
    <SearchInput source="rolesSearch" alwaysOn resettable placeholder="Roles" />
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
  <List
    pagination={<AdminPagination />}
    perPage={25}
    sort={defaultSort}
    filters={<RequestedOfficeUserListFilter />}
    actions={<ListActions />}
  >
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
