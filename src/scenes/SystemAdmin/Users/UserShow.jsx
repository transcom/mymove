import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField } from 'react-admin';

const UserShow = (props) => {
  return (
    <Show {...props} title="title">
      <SimpleShowLayout>
        <TextField source="loginGovEmail" label="user email" />
        <BooleanField source="active" />
        <TextField source="currentAdminSessionId" label="user current admin session ID" />
        <TextField source="currentOfficeSessionId" label="user current office session ID" />
        <TextField source="currentMilSessionId" label="user current mil session ID" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;
