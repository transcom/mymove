/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Datagrid, DateField, Filter, List, TextField, TextInput } from 'react-admin';
import { makeStyles } from '@material-ui/core/styles';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'locator', order: 'ASC' };

const MoveFilter = (props) => (
  <Filter {...props}>
    <TextInput label="Locator" source="locator" reference="locator" alwaysOn resettable />
  </Filter>
);

const useStyles = makeStyles({
  tableCell: {
    whiteSpace: 'normal',
    minWidth: '25px',
    padding: '8px',
  },
});

const MoveList = () => {
  const classes = useStyles();
  return (
    <List pagination={<AdminPagination />} perPage={25} filters={<MoveFilter />} sort={defaultSort}>
      <Datagrid bulkActionButtons={false} rowClick="show">
        <TextField
          source="id"
          reference="moves"
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <TextField
          source="ordersId"
          reference="moves"
          label="Order Id"
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <TextField
          source="serviceMember.id"
          label="Service Member Id"
          sortable={false}
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <TextField
          source="locator"
          reference="moves"
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <TextField
          source="status"
          reference="moves"
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <TextField
          source="show"
          reference="moves"
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <DateField
          source="createdAt"
          reference="moves"
          showTime
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <DateField
          source="updatedAt"
          reference="moves"
          showTime
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
        />
        <DateField
          source="availableToPrimeAt"
          reference="moves"
          showTime
          cellClassName={classes.tableCell}
          headerClassName={classes.tableCell}
          label="Available to Prime at"
        />
      </Datagrid>
    </List>
  );
};

export default MoveList;
