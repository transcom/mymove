import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField } from 'react-admin';

const UserShow = (props) => {
  return (
    <Show {...props} title="title" data-testid="user-show-detail">
      <SimpleShowLayout>
        <TextField data-testid="user-id" source="id" label="user ID" />
        <TextField data-testid="user-gov-email" source="loginGovEmail" label="user email" />
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
