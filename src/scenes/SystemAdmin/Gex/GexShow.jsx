import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField } from 'react-admin';

// const GexShowTitle = ({ record }) => {
//   return <span>{`${record.firstName} ${record.lastName}`}</span>;
// };

const GexShow = (props) => {
  return (
    // <Show {...props} title={<GexShowTitle />}>
    <Show {...props}>
      <h1>look at me... show</h1>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="lastName" />
        <TextField source="organizationId" />
        <BooleanField source="active" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default GexShow;
