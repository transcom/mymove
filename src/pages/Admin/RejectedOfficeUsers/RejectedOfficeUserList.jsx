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
} from 'react-admin';

import styles from './RejectedOfficeUserList.module.scss';

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

// Overriding the default toolbar
const ListActions = () => {
  return <TopToolbar />;
};

const filterList = [
  <SearchInput source="search" alwaysOn />,
  <TextInput label="Email" source="emails" />,
  <TextInput label="First Name" source="firstName" />,
  <TextInput label="Last Name" source="lastName" />,
  <TextInput label="Office" source="offices" />,
  <TextInput label="Rejection Reason" source="rejectionReason" />,
  <TextInput label="Rejected On" source="rejectedOn" />,
  <TextInput label="Roles" source="roles" />,
];

const SearchFilters = () => (
  <div className={styles.searchContainer}>
    <div className={styles.searchBar}>
      <FilterForm filters={filterList} />
    </div>
    <div className={styles.filters}>
      <FilterButton filters={filterList} />
    </div>
  </div>
);

const defaultSort = { field: 'createdAt', order: 'DESC' };

const RejectedOfficeUserList = () => (
  <List
    filters={<SearchFilters />}
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
        <RejectedOfficeUserShowRoles source="roles" />
      </ReferenceField>
    </Datagrid>
  </List>
);

export default RejectedOfficeUserList;
