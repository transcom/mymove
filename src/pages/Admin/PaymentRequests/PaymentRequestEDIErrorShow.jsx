import React from 'react';
import { Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';

const PaymentRequestEDIErrorShowTitle = false;
// const PaymentRequestEDIErrorShowTitle = () => {
//   const record = useRecordContext();
//   return <span>{`Payment Request EDI File Id: ${record.id}`}</span>;
// };

const CustomEdiStringField = ({ source }) => {
  const record = useRecordContext();
  return <div style={{ whiteSpace: 'pre-wrap', fontFamily: 'monospace' }}>{record[source]}</div>;
};

const PaymentRequestEDIErrorShow = () => {
  return (
    <Show title={<PaymentRequestEDIErrorShowTitle />}>
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

export default PaymentRequestEDIErrorShow;
