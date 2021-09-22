import React from 'react';
import { BooleanField, Datagrid, Filter, List, TextField, TextInput, TopToolbar } from 'react-admin';
import PropTypes from 'prop-types';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

// Overriding the default toolbar to add import button
const ListActions = () => {
  return <TopToolbar />;
};

const ClientCertListFilter = (props) => (
  <Filter {...props}>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'subject', order: 'ASC' };

const ClientCertList = (props) => (
  <List
    {...props}
    pagination={<AdminPagination />}
    perPage={25}
    bulkActionButtons={false}
    sort={defaultSort}
    filters={<ClientCertListFilter />}
    actions={<ListActions />}
  >
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="sha246Digest" />
      <TextField source="subject" />
      <BooleanField source="allowDpsAuthAPI" />
      <BooleanField source="allowOrdersAPI" />
      <BooleanField source="allowAirForceOrdersRead" />
      <BooleanField source="allowAirForceOrdersWrite" />
      <BooleanField source="allowArmyOrdersRead" />
      <BooleanField source="allowArmyOrdersWrite" />
      <BooleanField source="allowCoastGuardOrdersRead" />
      <BooleanField source="allowCoastGuardOrdersWrite" />
      <BooleanField source="allowMarineCorpsOrdersRead" />
      <BooleanField source="allowMarineCorpsOrdersWrite" />
      <BooleanField source="allowNavyOrdersRead" />
      <BooleanField source="allowNavyOrdersWrite" />
      <BooleanField source="allowPrime" />
    </Datagrid>
  </List>
);

ListActions.propTypes = {
  currentSort: PropTypes.exact({
    field: PropTypes.string,
    order: PropTypes.string,
  }),
  filterValues: PropTypes.shape({
    // This will have to be updated if we have any filters besides search added to this page
    search: PropTypes.string,
  }),
};

ListActions.defaultProps = {
  currentSort: {
    field: 'subject',
    order: 'ASC',
  },
  filterValues: {},
};

export default ClientCertList;
