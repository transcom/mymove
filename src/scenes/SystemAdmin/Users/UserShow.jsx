import React from 'react';
import { ShowController, ShowView, SimpleShowLayout, TextField } from 'react-admin';

const UserShow = (props) => (
  <ShowController {...props}>
    {(controllerProps) => (
      <ShowView {...props} {...controllerProps}>
        <SimpleShowLayout>
          <TextField source="user.loginGovEmail" label="user email" />
          <TextField source="user.currentAdminSessionId" label="user current admin session ID" />
          <TextField source="user.currentOfficeSessionId" label="user current office session ID" />
          <TextField source="user.currentMilSessionId" label="user current mil session ID" />
        </SimpleShowLayout>
      </ShowView>
    )}
  </ShowController>
);

export default UserShow;
