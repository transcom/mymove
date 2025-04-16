/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { SimpleShowLayout, Show, TextField } from 'react-admin';

const EdiErrorsShow = (props) => (
  <Show {...props}>
    <SimpleShowLayout>
      <TextField source="id" />
      <TextField source="paymentRequestID" />
      <TextField source="paymentRequestNumber" />
      <TextField source="ediType" />
      <TextField source="code" />
      <TextField source="description" />
      <TextField source="createdAt" />
    </SimpleShowLayout>
  </Show>
);

export default EdiErrorsShow;
