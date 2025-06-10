import React from 'react';
import {
  Datagrid,
  DateField,
  List,
  ReferenceField,
  TextField,
  TextInput,
  TopToolbar,
  useRecordContext,
  SearchInput,
  FilterForm,
  FilterButton,
  downloadCSV,
  useListController,
  ExportButton,
} from 'react-admin';
import * as jsonexport from 'jsonexport/dist';

import styles from './RejectedOfficeUserList.module.scss';

import { getTransportationOfficeByID } from 'services/adminApi';
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

  return <span>{uniqueRoleNamesList.join(', ')}</span>;
};

// Custom exporter to flatten out role  types
const exporter = async (data) => {
  const usersForExportPromises = data.map(async (rowData) => {
    const { roles, email, firstName, lastName, status, rejectionReason, rejectedOn, transportationOfficeId } = rowData;

    const flattenedRoles = roles ? [...new Set(roles.map((role) => role.roleName))].join(',') : '';
    const transportationOffice = await getTransportationOfficeByID(transportationOfficeId);
    const officeName = transportationOffice.name;
    return {
      email,
      firstName,
      lastName,
      status,
      rejectionReason,
      rejectedOn,
      officeName,
      roles: flattenedRoles,
    };
  });

  const usersForExport = await Promise.all(usersForExportPromises);

  // convert data to csv and download
  jsonexport(usersForExport, {}, (err, csv) => {
    if (err) throw err;
    downloadCSV(csv, 'rejected-office-users');
  });
};

const filterList = [
  <SearchInput source="search" alwaysOn />,
  <TextInput label="Email" source="emails" />,
  <TextInput label="First Name" source="firstName" />,
  <TextInput label="Last Name" source="lastName" />,
  <TextInput label="Office" source="offices" />,
  <TextInput label="Rejection Reason" source="rejectionReason" />,
  <TextInput label="Rejected On" placeholder="MM/DD/YYYY" source="rejectedOn" />,
  <TextInput label="Roles" source="roles" />,
];

const RejectedOfficeUserListFilter = () => (
  <div className={styles.searchContainer}>
    <div className={styles.searchBar}>
      <FilterForm filters={filterList} />
    </div>
  </div>
);

// Overriding the default toolbar
const ListActions = () => {
  const { total, resource, sort, filterValues } = useListController();

  return (
    <TopToolbar>
      <FilterButton filters={filterList} />
      <ExportButton disabled={total === 0} resource={resource} sort={sort} filter={filterValues} exporter={exporter} />
    </TopToolbar>
  );
};

const defaultSort = { field: 'createdAt', order: 'DESC' };

const RejectedOfficeUserList = () => (
  <List
    filters={<RejectedOfficeUserListFilter />}
    pagination={<AdminPagination />}
    perPage={25}
    sort={defaultSort}
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
      <DateField showTime source="rejectedOn" label="Rejected on" />
      <ReferenceField label="Roles Requested" source="id" sortBy="role" reference="rejected-office-users" link={false}>
        <RejectedOfficeUserShowRoles sortable={false} source="roles" label="Rejected Roles" />
      </ReferenceField>
    </Datagrid>
  </List>
);

export default RejectedOfficeUserList;
