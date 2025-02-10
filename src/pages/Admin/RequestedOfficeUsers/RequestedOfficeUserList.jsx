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
  TopToolbar,
  useRecordContext,
  BulkExportButton,
  downloadCSV,
  useListContext,
  useDataProvider,
} from 'react-admin';
import jsonExport from 'jsonexport/dist';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar
const ListActions = () => {
  return <TopToolbar />;
};
const RequestedOfficeUserListFilter = (props) => (
  <Filter {...props}>
    <TextInput label="Search by Name/Email" source="search" alwaysOn resettable />
    <TextInput label="Transportation Office" source="offices" alwaysOn resettable />
    <TextInput label="Roles" source="rolesSearch" alwaysOn resettable />
  </Filter>
);

const defaultSort = { field: 'createdAt', order: 'DESC' };

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

const handleExport = async (data, dataProvider, selectedIds) => {
  const allRequestedOfficeUsers = data;

  const selectedUserIdObjects = {};
  selectedIds.forEach((id) => {
    if (!selectedUserIdObjects[`${id}`]) {
      selectedUserIdObjects[`${id}`] = id;
    }
  });

  const selectedUsers = [];
  allRequestedOfficeUsers.forEach((user) => {
    if (selectedUserIdObjects[`${user.id}`]) {
      selectedUsers.push(user);
    }
  });

  const officeObjects = {};
  const offices = await dataProvider.getMany('offices');
  offices.data.forEach((office) => {
    if (!officeObjects[`${office.id}`]) {
      officeObjects[`${office.id}`] = office;
    }
  });

  const usersWithTransportationOfficeName = selectedUsers.map((user) => ({
    ...user,
    transportationOfficeName: officeObjects[user.transportationOfficeId]?.name,
  }));

  const userToCSV = [];
  usersWithTransportationOfficeName.forEach((user) => {
    const userRoles = UserRolesToString(user);
    const csvUser = {
      Id: user.id,
      Email: user.email,
      'First Name': user.firstName,
      'Last Name': user.lastName,
      'Transportation Office': user.transportationOfficeName,
      Status: user.status,
      'Requested On': user.createdAt,
      Roles: userRoles,
    };

    userToCSV.push(csvUser);
  });

  jsonExport(userToCSV, (err, csv) => {
    downloadCSV(csv, 'Requested_Office_Users');
  });
};

const CustomBulkActions = ({ selectedIds }) => {
  const { data } = useListContext();
  const dataProvider = useDataProvider();

  return <BulkExportButton label="Export" onClick={() => handleExport(data, dataProvider, selectedIds)} />;
};

const RequestedOfficeUserList = () => {
  return (
    <List
      pagination={<AdminPagination />}
      perPage={25}
      sort={defaultSort}
      filters={<RequestedOfficeUserListFilter />}
      actions={<ListActions />}
    >
      <Datagrid bulkActionButtons={<CustomBulkActions />} rowClick="show" data-testid="requested-office-user-fields">
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
