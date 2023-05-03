import React from 'react';
import { BooleanField, DateField, Show, SimpleShowLayout, TextField, useRecordContext } from 'react-admin';

const AdminUserShowTitle = () => {
  const record = useRecordContext();
  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const AdminUserShow = () => {
  return (
    <Show title={<AdminUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="lastName" />
        <TextField source="organizationId" label="Organization Id" />
        <BooleanField source="active" label="Active" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default AdminUserShow;
