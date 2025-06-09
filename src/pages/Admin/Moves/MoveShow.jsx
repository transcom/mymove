import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';

const MoveShowTitle = () => {
  const record = useRecordContext();
  return <span>{`Move ID: ${record.id}`}</span>;
};

const MoveShow = () => {
  return (
    <Show title={<MoveShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="locator" />
        <TextField source="status" />
        <BooleanField source="show" />
        <TextField source="ordersId" reference="moves" label="Order Id" />
        <TextField source="serviceMember.userId" label="User Id" />
        <TextField source="serviceMember.id" label="Service member Id" />
        <TextField source="serviceMember.firstName" label="Service member first name" />
        <TextField source="serviceMember.middleName" label="Service member middle name" />
        <TextField source="serviceMember.lastName" label="Service member last name" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
        <DateField source="availableToPrimeAt" showTime label="Available to Prime at" />
      </SimpleShowLayout>
    </Show>
  );
};

export default MoveShow;
