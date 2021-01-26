import React from 'react';
import { Show, SimpleShowLayout, TextField } from 'react-admin';

const UserShow = (props) => {
  return (
    <Show {...props} title="title">
      <SimpleShowLayout>
        <TextField source="loginGovEmail" label="user email" />
        <TextField source="currentAdminSessionId" label="user current admin session ID" />
        <TextField source="currentOfficeSessionId" label="user current office session ID" />
        <TextField source="currentMilSessionId" label="user current mil session ID" />
      </SimpleShowLayout>
    </Show>
  );
};

export default UserShow;
