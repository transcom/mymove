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
} from 'react-admin';
import PropTypes from 'prop-types';

import ImportOfficeUserButton from 'components/Admin/ImportOfficeUserButton';
import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar to add import button
const ListActions = (props) => {
  const { total, resource, sort, filterValues, exporter } = props;
  return (
    <TopToolbar>
      <CreateButton />
      <ImportOfficeUserButton resource={resource} {...props} />
      <ExportButton disabled={total === 0} resource={resource} sort={sort} filter={filterValues} exporter={exporter} />
    </TopToolbar>
  );
};

const OfficeUserListFilter = (props) => (
  <Filter {...props}>
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
  resource: PropTypes.string,
  sort: PropTypes.exact({
    field: PropTypes.string,
    order: PropTypes.string,
  }),
  filterValues: PropTypes.shape({
    // This will have to be updated if we have any filters besides search added to this page
    search: PropTypes.string,
  }),
  exporter: PropTypes.func.isRequired,
};

ListActions.defaultProps = {
  resource: 'office_users',
  sort: {
    field: 'last_name',
    order: 'ASC',
  },
  basePath: undefined,
  total: null,
  filterValues: {},
};

export default OfficeUserList;
