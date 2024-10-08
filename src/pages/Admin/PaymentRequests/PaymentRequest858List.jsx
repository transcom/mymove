/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Datagrid, Filter, List, TextField, TextInput } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'createdAt', order: 'ASC' };

const PaymentRequestFilter = (props) => (
  <Filter {...props}>
    <TextInput
      label="Payment Request Number"
      source="paymentRequestNumber"
      reference="paymentRequestNumber"
      alwaysOn
      resettable
    />
  </Filter>
);

const PaymentRequest858List = () => (
  <List pagination={<AdminPagination />} filters={PaymentRequestFilter} perPage={25} sort={defaultSort}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" />
      <TextField source="paymentRequestNumber" />
      <TextField source="fileName" />
      <TextField source="ediString" />
      <TextField source="createdAt" />
    </Datagrid>
  </List>
);

export default PaymentRequest858List;
