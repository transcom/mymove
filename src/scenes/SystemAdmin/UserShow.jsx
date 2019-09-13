import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField } from 'react-admin';

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
        <TextField source="last_name" />
        <BooleanField source="disabled" label="Deactivated" />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;
