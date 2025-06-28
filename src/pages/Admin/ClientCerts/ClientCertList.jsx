import React from 'react';
import { BooleanField, Datagrid, Filter, List, TextField, TextInput } from 'react-admin';
import { makeStyles } from '@material-ui/core/styles';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const ClientCertListFilter = (props) => (
  // eslint-disable-next-line react/jsx-props-no-spreading
  <Filter {...props}>
    <TextInput labl="Search by Cert Subject" source="search" resettable alwaysOn />
  </Filter>
);

const defaultSort = { field: 'subject', order: 'ASC' };

const useStyles = makeStyles(() => ({
  tableCell: {
    minWidth: 25,
    maxWidth: 150,
    whiteSpace: 'normal',
    overflow: 'scroll',
    overflowWrap: 'break-word',
  },
}));

const useHeaderStyles = makeStyles(() => ({
  headerCell: {
    minWidth: 25,
  },
}));

const ClientCertList = (props) => {
  const classes = useStyles();
  const headerClasses = useHeaderStyles();
  return (
    <List
      // eslint-disable-next-line react/jsx-props-no-spreading
      {...props}
      pagination={<AdminPagination />}
      perPage={25}
      bulkActionButtons={false}
      sort={defaultSort}
      filters={<ClientCertListFilter />}
    >
      <Datagrid rowClick="show" classes={headerClasses}>
        <TextField cellClassName={classes.tableCell} source="subject" />
        <TextField source="id" />
        <TextField source="sha256Digest" />
        <TextField source="userId" label="User Id" />
        <BooleanField cellClassName={classes.tableCell} source="allowPrime" label="Prime API" />
        <BooleanField cellClassName={classes.tableCell} source="allowPPTAS" label="PPTAS API" />
        <TextField cellClassName={classes.tableCell} source="pptasAffiliation" label="PPTAS Affiliation" />
        <BooleanField cellClassName={classes.tableCell} source="allowOrdersAPI" label="Orders API" />
        <BooleanField cellClassName={classes.tableCell} source="allowAirForceOrdersRead" label="USAF Orders Read" />
        <BooleanField cellClassName={classes.tableCell} source="allowAirForceOrdersWrite" label="USAF Orders Write" />
        <BooleanField cellClassName={classes.tableCell} source="allowArmyOrdersRead" label="Army Orders Read" />
        <BooleanField cellClassName={classes.tableCell} source="allowArmyOrdersWrite" label="Army Orders Write" />
        <BooleanField cellClassName={classes.tableCell} source="allowCoastGuardOrdersRead" label="USCG Orders Read" />
        <BooleanField cellClassName={classes.tableCell} source="allowCoastGuardOrdersWrite" label="USCG Orders Write" />
        <BooleanField cellClassName={classes.tableCell} source="allowMarineCorpsOrdersRead" label="USMC Orders Read" />
        <BooleanField
          cellClassName={classes.tableCell}
          source="allowMarineCorpsOrdersWrite"
          label="USMC Orders Write"
        />
        <BooleanField cellClassName={classes.tableCell} source="allowNavyOrdersRead" label="Navy Orders Read" />
        <BooleanField cellClassName={classes.tableCell} source="allowNavyOrdersWrite" label="Navy Orders Write" />
      </Datagrid>
    </List>
  );
};

export default ClientCertList;
