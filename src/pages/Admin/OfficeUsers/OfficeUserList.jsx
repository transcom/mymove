import React from 'react';
import {
  BooleanField,
  CreateButton,
  Datagrid,
  ExportButton,
  Filter,
  List,
  ReferenceField,
  TextField,
  TextInput,
  TopToolbar,
  useListController,
} from 'react-admin';

import ImportOfficeUserButton from 'components/Admin/ImportOfficeUserButton';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar to add import button
const ListActions = () => {
  const { total, resource, sort, filterValues, exporter } = useListController();
  return (
    <TopToolbar>
      <CreateButton />
      <ImportOfficeUserButton resource={resource} />
      <ExportButton disabled={total === 0} resource={resource} sort={sort} filter={filterValues} exporter={exporter} />
    </TopToolbar>
  );
};

const OfficeUserListFilter = () => (
  <Filter>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'last_name', order: 'ASC' };

const OfficeUserList = () => (
  <List
    pagination={<AdminPagination />}
    perPage={25}
    sort={defaultSort}
    filters={<OfficeUserListFilter />}
    actions={<ListActions />}
  >
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="firstName" />
      <TextField source="lastName" />
      <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" link={false}>
        <TextField source="name" />
      </ReferenceField>
      <TextField source="userId" label="User Id" />
      <BooleanField source="active" />
    </Datagrid>
  </List>
);

export default OfficeUserList;
