/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Datagrid, Filter, List, TextField, TextInput, TopToolbar } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const defaultSort = { field: 'createdAt', order: 'ASC' };

const ListActions = () => {
  return <TopToolbar />;
};

const PaymentRequestFilter = (props) => (
  <Filter {...props}>
    <TextInput source="paymentRequestNumber" alwaysOn resettable />
  </Filter>
);

const PaymentRequest858List = () => (
  <List
    pagination={<AdminPagination />}
    filters={<PaymentRequestFilter />}
    perPage={25}
    sort={defaultSort}
    actions={<ListActions />}
  >
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="id" />
      <TextField source="paymentRequestNumber" />
      <TextField source="fileName" />
      <TextField source="createdAt" />
    </Datagrid>
  </List>
);

export default PaymentRequest858List;
