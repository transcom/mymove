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
  downloadCSV,
  useListController,
  ExportButton,
} from 'react-admin';
import * as jsonexport from 'jsonexport/dist';

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

// Custom exporter to flatten out role  types
const exporter = (data) => {
  const usersForExport = data.map((rowData) => {
    const { roles, email, firstName, lastName, status, rejectionReason, rejectedOn } = rowData;

    const flattenedRoles = roles ? [...new Set(roles.map((role) => role.roleName))].join(',') : '';

    return {
      email,
      firstName,
      lastName,
      status,
      rejectionReason,
      rejectedOn,

      roles: flattenedRoles,
    };
  });

  // convert data to csv and download
  jsonexport(usersForExport, {}, (err, csv) => {
    if (err) throw err;
    downloadCSV(csv, 'rejected-office-users');
  });
};

// Overriding the default toolbar
const ListActions = () => {
  // return <TopToolbar />;
  const { total, resource, sort, filterValues } = useListController();

  return (
    <TopToolbar>
      <ExportButton disabled={total === 0} resource={resource} sort={sort} filter={filterValues} exporter={exporter} />
    </TopToolbar>
  );
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
