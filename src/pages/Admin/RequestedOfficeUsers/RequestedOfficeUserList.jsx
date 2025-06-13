import React from 'react';
import {
  ArrayField,
  Datagrid,
  DateField,
  SearchInput,
  FilterForm,
  FilterButton,
  List,
  ReferenceField,
  TextField,
  TextInput,
  TopToolbar,
  useRecordContext,
  downloadCSV,
  useDataProvider,
  ExportButton,
  useListController,
} from 'react-admin';
import jsonExport from 'jsonexport/dist';

import styles from './RequestedOfficeUserList.module.scss';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const UserRolesToString = (user) => {
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
  return <div>{UserRolesToString(record)}</div>;
};

const RequestedOfficeUserListFilter = [
  <SearchInput source="search" alwaysOn />,
  <TextInput label="Email" source="email" />,
  <TextInput label="First Name" source="firstName" />,
  <TextInput label="Last Name" source="lastName" />,
  <TextInput label="Office" source="office" />,
  <TextInput label="Requested On" placeholder="MM/DD/YYYY" source="requestedOn" />,
  <TextInput label="Roles" source="roles" />,
];

const ListActions = () => {
  const { total, resource, sort, filterValues } = useListController();
  const dataProvider = useDataProvider();

  const exporter = async (users) => {
    const officeObjects = {};
    const offices = await dataProvider.getMany('offices');
    offices.data.forEach((office) => {
      if (!officeObjects[`${office.id}`]) {
        officeObjects[`${office.id}`] = office;
      }
    });

    const usersWithTransportationOfficeName = users.map((user) => ({
      ...user,
      transportationOfficeName: officeObjects[user.transportationOfficeId]?.name,
    }));

    const usersForExport = usersWithTransportationOfficeName.map((user) => {
      const { id, email, firstName, lastName, transportationOfficeName, status, createdAt } = user;
      const userRoles = UserRolesToString(user);
      return {
        Id: id,
        Email: email,
        'First Name': firstName,
        'Last Name': lastName,
        'Transportation Office': transportationOfficeName,
        Status: status,
        'Requested On': createdAt,
        Roles: userRoles,
      };
    });

    // convert data to csv and download
    jsonExport(usersForExport, {}, (err, csv) => {
      if (err) throw err;
      downloadCSV(csv, 'requested_office_users');
    });
  };

  return (
    <TopToolbar>
      <FilterButton filters={RequestedOfficeUserListFilter} />
      <ExportButton disabled={total === 0} resource={resource} sort={sort} filter={filterValues} exporter={exporter} />
    </TopToolbar>
  );
};

const SearchFilters = () => (
  <div className={styles.searchContainer}>
    <div className={styles.searchBar}>
      <FilterForm filters={RequestedOfficeUserListFilter} />
    </div>
  </div>
);

const defaultSort = { field: 'createdAt', order: 'DESC' };

const RequestedOfficeUserList = () => {
  return (
    <List
      pagination={<AdminPagination />}
      perPage={25}
      sort={defaultSort}
      filters={<SearchFilters />}
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
};

export default RequestedOfficeUserList;
