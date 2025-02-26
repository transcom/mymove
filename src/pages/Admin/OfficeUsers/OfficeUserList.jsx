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
  useDataProvider,
} from 'react-admin';
import * as jsonexport from 'jsonexport/dist';

import { ADMIN_APP_USER_EXPORT_HEADERS } from '../../../shared/constants';

import ImportOfficeUserButton from 'components/Admin/ImportOfficeUserButton';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Function to transform rowData based on headers
const transformRowData = (rowData, officeObjects) => {
  const transformedData = {};
  ADMIN_APP_USER_EXPORT_HEADERS.forEach(({ key, header }) => {
    switch (key) {
      case 'roles':
        transformedData[header] = rowData[key] ? rowData[key].map((role) => role.roleType).join(',') : '';
        break;
      case 'privileges':
        transformedData[header] = rowData[key]
          ? rowData[key].map((privilege) => privilege.privilegeType).join(',')
          : '';
        break;
      case 'primaryTransportationOffice':
        transformedData[header] = officeObjects[rowData.transportationOfficeId]?.name || '';
        break;
      default:
        transformedData[header] = rowData[key] !== undefined ? rowData[key] : '';
        break;
    }
  });
  return transformedData;
};

// Overriding the default toolbar for customizations
const ListActions = () => {
  const { total, resource, sort, filterValues } = useListController();
  const dataProvider = useDataProvider();
  const exporter = async (data) => {
    // Fetch the offices asynchronously
    const officesResponse = await dataProvider.getMany('offices');
    const officeObjects = {};

    // Map office data into officeObjects using the office id as the key
    officesResponse.data.forEach((office) => {
      if (!officeObjects[`${office.id}`]) {
        officeObjects[`${office.id}`] = office;
      }
    });

    // Process the user data using the transformation function
    const usersForExport = data.map((rowData) => transformRowData(rowData, officeObjects));

    // Extract header names for jsonexport
    const headersMap = ADMIN_APP_USER_EXPORT_HEADERS.map((h) => h.header);
    // Convert the data to CSV and trigger the download
    jsonexport(usersForExport, { headersMap }, (err, csv) => {
      if (err) throw err;
      downloadCSV(csv, 'office-users');
    });
  };

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
      <TextField source="telephone" />
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
