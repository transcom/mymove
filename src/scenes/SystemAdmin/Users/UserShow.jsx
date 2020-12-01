import React from 'react';
import { ShowController, ShowView, SimpleShowLayout, TextField } from 'react-admin';

const UserShow = (props) => (
  <ShowController {...props}>
    {(controllerProps) => (
      <ShowView {...props} {...controllerProps}>
        <SimpleShowLayout>
          <TextField source="loginGovEmail" label="user email" />
          <TextField source="currentAdminSessionId" label="user current admin session ID" />
          <TextField source="currentOfficeSessionId" label="user current office session ID" />
          <TextField source="currentMilSessionId" label="user current mil session ID" />
        </SimpleShowLayout>
      </ShowView>
    )}
  </ShowController>
);

export default UserShow;
