import React from 'react';
import { Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';

const PaymentRequest858ShowTitle = () => {
  const record = useRecordContext();
  return <span>{`Payment Request EDI File Id: ${record.id}`}</span>;
};

const CustomEdiStringField = ({ source }) => {
  const record = useRecordContext();
  return <div style={{ whiteSpace: 'pre-wrap', fontFamily: 'monospace' }}>{record[source]}</div>;
};

const PaymentRequest858Show = () => {
  return (
    <Show title={<PaymentRequest858ShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="paymentRequestNumber" />
        <TextField source="fileName" />
        <CustomEdiStringField source="ediString" />
        <TextField source="createdAt" />
      </SimpleShowLayout>
    </Show>
  );
};

export default PaymentRequest858Show;
