import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField } from 'react-admin';

const UserShowTitle = ({ record }) => {
  return <span>{`${record.first_name} ${record.last_name}`}</span>;
};

const UserShow = props => {
  return (
    <Show {...props} title={<UserShowTitle />}>
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="email" />
        <TextField source="first_name" />
        <TextField source="middle_initials" />
        <TextField source="last_name" />
        <BooleanField source="disabled" label="Deactivated" />
        <DateField source="created_at" showTime />
        <DateField source="updated_at" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;
