import React from 'react';
import {
  BooleanField,
  Datagrid,
  Filter,
  List,
  ReferenceField,
  TextField,
  TextInput,
  TopToolbar,
  CreateButton,
  ExportButton,
} from 'react-admin';
import { ImportButton } from 'react-admin-import-csv';
import PropTypes from 'prop-types';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const avaliableRoles = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
  { roleType: 'transportation_invoicing_officer', name: 'Transportation Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'ppm_office_users', name: 'PPM Office Users' },
  { roleType: 'services_counselor', name: 'Services Counselor' },
];

const validateRow = async () => {};

const preCommitCallback = (action, rows) => {
  const alteredRows = [];
  rows.forEach((row) => {
    const copyOfRow = row;
    if (row.roles) {
      const rolesArray = [];

      // Parse roles from string
      const parsedRoles = row.roles.split(',');
      parsedRoles.forEach((parsedRole) => {
        // Remove any whitespace in the role string
        const role = parsedRole.replaceAll(/\s/g, '');
        rolesArray.push(avaliableRoles.find((avaliableRole) => avaliableRole.roleType === role));
      });
      copyOfRow.roles = rolesArray;
    }
    alteredRows.push(copyOfRow);
  });
  return alteredRows;
};

const config = {
  logging: true,
  validateRow,
  preCommitCallback,
  postCommitCallback: () => {
    // console.log('reportItems', { reportItems });
  },
  disableImportOverwrite: true,
};

// Overriding the default toolbar to add import button
const ListActions = (props) => {
  const { basePath, total, resource, currentSort, filterValues, exporter } = props;
  return (
    <TopToolbar>
      <CreateButton basePath={basePath} />
      <ImportButton {...props} {...config} />
      <ExportButton
        disabled={total === 0}
        resource={resource}
        sort={currentSort}
        filter={filterValues}
        exporter={exporter}
      />
    </TopToolbar>
  );
};

const OfficeUserListFilter = (props) => (
  <Filter {...props}>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'last_name', order: 'ASC' };

const OfficeUserList = (props) => (
  <List
    {...props}
    pagination={<AdminPagination />}
    perPage={25}
    bulkActionButtons={false}
    sort={defaultSort}
    filters={<OfficeUserListFilter />}
    actions={<ListActions />}
  >
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="email" />
      <TextField source="firstName" />
      <TextField source="lastName" />
      <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices">
        <TextField source="name" />
      </ReferenceField>
      <TextField source="userId" label="User Id" />
      <BooleanField source="active" />
    </Datagrid>
  </List>
);

ListActions.propTypes = {
  basePath: PropTypes.string,
  total: PropTypes.number,
  resource: PropTypes.string.isRequired,
  currentSort: PropTypes.exact({
    field: PropTypes.string,
    order: PropTypes.string,
  }).isRequired,
  filterValues: PropTypes.shape({
    // This will have to be updated if we have any filters besides search added to this page
    search: PropTypes.string,
  }),
  exporter: PropTypes.func.isRequired,
};

ListActions.defaultProps = {
  basePath: undefined,
  total: null,
  filterValues: {},
};

export default OfficeUserList;
