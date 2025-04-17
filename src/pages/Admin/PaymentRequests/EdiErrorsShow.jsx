/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { SimpleShowLayout, Show, TextField, DateField } from 'react-admin';

const EdiErrorsShow = (props) => (
  <Show {...props}>
    <SimpleShowLayout>
      <TextField source="id" label="EDI Error ID" />
      <TextField source="paymentRequestID" label="Payment Request ID" />
      <TextField source="paymentRequestNumber" label="Payment Request Number" />
      <TextField source="ediType" label="Error Type" />
      <TextField source="code" label="Error Code" />
      <TextField source="description" label="Error Description" />
      <DateField source="createdAt" showTime label="Error Created At" sortable={false} />
    </SimpleShowLayout>
  </Show>
);

export default EdiErrorsShow;
