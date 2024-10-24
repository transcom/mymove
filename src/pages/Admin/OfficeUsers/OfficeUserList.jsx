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
  downloadCSV,
} from 'react-admin';
import * as jsonexport from 'jsonexport/dist';

import ImportOfficeUserButton from 'components/Admin/ImportOfficeUserButton';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Custom exporter to flatten out role and privilege types
const exporter = (data) => {
  const usersForExport = data.map((rowData) => {
    const { roles, privileges, ...otherRowData } = rowData;

    const flattenedRoles = roles ? roles.map((role) => role.roleType).join(',') : '';
    const flattenedPrivileges = privileges ? privileges.map((privilege) => privilege.privilegeType).join(',') : '';

    return {
      ...otherRowData,
      roles: flattenedRoles,
      privileges: flattenedPrivileges,
    };
  });

  // convert data to csv and download
  jsonexport(usersForExport, {}, (err, csv) => {
    if (err) throw err;
    downloadCSV(csv, 'office-users');
  });
};

// Overriding the default toolbar to add import button
const ListActions = () => {
  const { total, resource, sort, filterValues } = useListController();

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
      <ReferenceField
        label="Primary Transportation Office"
        source="transportationOfficeId"
        reference="offices"
        link={false}
      >
        <TextField source="name" />
      </ReferenceField>
      <TextField source="userId" label="User Id" />
      <BooleanField source="active" />
    </Datagrid>
  </List>
);

export default OfficeUserList;
