import React from 'react';
import { BooleanField, Datagrid, Filter, List, TextField, TextInput } from 'react-admin';
import { makeStyles } from '@material-ui/core/styles';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const ClientCertListFilter = (props) => (
  <Filter {...props}>
    <TextInput source="search" alwaysOn />
  </Filter>
);

const defaultSort = { field: 'subject', order: 'ASC' };

const useStyles = makeStyles(() => ({
  tableCell: {
    maxWidth: 150,
    whiteSpace: 'normal',
    overflow: 'scroll',
    overflowWrap: 'break-word',
  },
}));

const ClientCertList = ({ ...props }) => {
  const classes = useStyles();
  return (
    <List
      {...props}
      pagination={<AdminPagination />}
      perPage={25}
      bulkActionButtons={false}
      sort={defaultSort}
      filters={<ClientCertListFilter />}
    >
      <Datagrid rowClick="show">
        <TextField cellClassName={classes.tableCell} source="subject" />
        <BooleanField source="allowDpsAuthAPI" label="Allow DPS Auth API" />
        <BooleanField source="allowOrdersAPI" label="Allow Orders API" />
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
        <TextField source="id" />
        <TextField source="sha256Digest" />
      </Datagrid>
    </List>
  );
};

export default ClientCertList;
