import React from 'react';
import { Datagrid, DateField, List, TextField } from 'react-admin';

import AdminPagination from 'scenes/SystemAdmin/shared/AdminPagination';

const EdiErrorsList = () => (
  <List pagination={<AdminPagination />} perPage={25}>
    <Datagrid bulkActionButtons={false} rowClick="show">
      <TextField source="paymentRequestNumber" label="Payment Request Number" />
      <TextField source="code" label="Error Code" />
      <TextField source="ediType" label="Error Type" />
      <TextField source="description" label="Error Description" />
      <DateField source="createdAt" showTime label="Error Created At" sortable={false} />
    </Datagrid>
  </List>
);

export default EdiErrorsList;
