import React from 'react';
import { Show, SimpleShowLayout, TextField, BooleanField, DateField } from 'react-admin';

const UserShow = (props) => {
  return (
    // eslint-disable-next-line
    <Show {...props} title="title" data-testid="user-show-detail">
      <SimpleShowLayout>
        <TextField source="id" label="User ID" />
        <TextField source="loginGovEmail" label="User email" />
        <BooleanField source="active" />
        <TextField source="currentAdminSessionId" label="User current admin session ID" />
        <TextField source="currentOfficeSessionId" label="User current office session ID" />
        <TextField source="currentMilSessionId" label="User current mil session ID" />
        <DateField source="createdAt" showTime />
        <DateField source="updatedAt" showTime />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;
