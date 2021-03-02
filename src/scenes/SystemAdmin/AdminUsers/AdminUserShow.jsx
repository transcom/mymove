import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField } from 'react-admin';

const AdminUserShowTitle = ({ record }) => {
  return <span>{`${record.firstName} ${record.lastName}`}</span>;
};

const AdminUserShow = (props) => {
  return (
    <Show {...props} title={<AdminUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
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

export default AdminUserShow;
