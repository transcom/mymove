import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField } from 'react-admin';

const AdminUserShowTitle = ({ record }) => {
  return <span>{`${record.first_name} ${record.last_name}`}</span>;
};

const AdminUserShow = props => {
  return (
    <Show {...props} title={<AdminUserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="email" />
        <TextField source="first_name" />
        <TextField source="last_name" />
        <TextField source="organization_id" />
        <BooleanField source="active" />
        <DateField source="created_at" showTime />
        <DateField source="updated_at" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default AdminUserShow;
