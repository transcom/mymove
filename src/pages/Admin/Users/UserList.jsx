import React from 'react';
import { List, Datagrid, TextField, BooleanField, Filter, TextInput } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const UserFilter = (props) => (
  // eslint-disable-next-line react/jsx-props-no-spreading
  <Filter {...props}>
    <TextInput label="Search by User Id or Email" source="search" resettable alwaysOn />
  </Filter>
);

const defaultSort = { field: 'oktaEmail', order: 'ASC' };

const UserList = () => (
  <List
    filters={<UserFilter />}
    pagination={<AdminPagination />}
    perPage={25}
    sort={defaultSort}
    bulkActionButtons={false}
  >
    <Datagrid bulkActionButtons={false} data-testid="user-data-grid" rowClick="show">
      <TextField data-testid="user-id" source="id" />
      <TextField source="oktaEmail" label="Email" />
      <BooleanField source="active" />
      <TextField source="createdAt" />
    </Datagrid>
  </List>
);

export default UserList;
