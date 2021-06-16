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
import PropTypes from 'prop-types';

import ImportCsvButton from 'components/Admin/ImportCsvButton';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar to add import button
const ListActions = (props) => {
  const { basePath, total, resource, currentSort, filterValues, exporter } = props;
  return (
    <TopToolbar>
      <CreateButton basePath={basePath} />
      <ImportCsvButton {...props} />
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
